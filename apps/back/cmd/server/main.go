// Package main is the entrypoint for the seonology-journey-back gRPC service.
package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/seonNoh/seonology-journey/apps/back/internal/accommodation"
	"github.com/seonNoh/seonology-journey/apps/back/internal/day"
	"github.com/seonNoh/seonology-journey/apps/back/internal/meal"
	"github.com/seonNoh/seonology-journey/apps/back/internal/media"
	"github.com/seonNoh/seonology-journey/apps/back/internal/observability"
	"github.com/seonNoh/seonology-journey/apps/back/internal/record"
	"github.com/seonNoh/seonology-journey/apps/back/internal/repository/ddb"
	"github.com/seonNoh/seonology-journey/apps/back/internal/schedule"
	"github.com/seonNoh/seonology-journey/apps/back/internal/server"
	"github.com/seonNoh/seonology-journey/apps/back/internal/social"
	"github.com/seonNoh/seonology-journey/apps/back/internal/trip"
	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger := observability.NewLogger()
	slog.SetDefault(logger)

	// Metrics
	reg := prometheus.DefaultRegisterer
	metrics := observability.NewMetrics(reg)
	_ = metrics // used by interceptors (TODO: wire gRPC interceptor)

	// Health checker
	hc := observability.NewHealthChecker()
	hc.Set("grpc", observability.StatusUp)

	// Observability HTTP server (/metrics + /healthz)
	obsAddr := os.Getenv("OBS_LISTEN_ADDR")
	if obsAddr == "" {
		obsAddr = ":9091"
	}
	obsMux := http.NewServeMux()
	obsMux.Handle("/metrics", observability.Handler())
	obsMux.Handle("/healthz", hc.Handler())
	obsServer := &http.Server{Addr: obsAddr, Handler: obsMux}
	go func() {
		logger.Info("observability HTTP server starting", "addr", obsAddr)
		if err := obsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("observability server failed", "error", err)
		}
	}()

	addr := os.Getenv("GRPC_LISTEN_ADDR")
	if addr == "" {
		addr = ":9090"
	}
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("listen failed", "error", err)
		os.Exit(1)
	}

	// DynamoDB client
	ddbClient := ddb.Client(context.Background())

	// Verify required DDB tables exist at startup
	requiredTables := []string{
		"seonology-journey-trips",
		"seonology-journey-days",
		"seonology-journey-schedules",
		"seonology-journey-media",
		"seonology-journey-records",
		"seonology-journey-social",
	}
	for _, table := range requiredTables {
		_, err := ddbClient.DescribeTable(context.Background(), &dynamodb.DescribeTableInput{
			TableName: &table,
		})
		if err != nil {
			logger.Warn("DDB table may not exist", "table", table, "error", err)
		}
	}

	// Trip: DynamoDB repository
	tripRepo := trip.NewDDBRepo(ddbClient)

	// Day / Schedule / Meal / Accommodation: DynamoDB repositories
	dayRepo := day.NewDDBRepo(ddbClient)
	scheduleRepo := schedule.NewDDBRepo(ddbClient)
	mealRepo := meal.NewDDBRepo(ddbClient)
	accomRepo := accommodation.NewDDBRepo(ddbClient)

	expRepo := record.NewExpenseDDBRepo(ddbClient)
	noteRepo := record.NewNoteDDBRepo(ddbClient)
	chkRepo := record.NewChecklistDDBRepo(ddbClient)
	rsvRepo := record.NewReservationDDBRepo(ddbClient)

	compRepo := social.NewCompanionDDBRepo(ddbClient)
	tagRepo := social.NewTagDDBRepo(ddbClient)
	tplRepo := social.NewTemplateDDBRepo(ddbClient)
	favRepo := social.NewFavoriteDDBRepo(ddbClient)
	shareRepo := social.NewShareDDBRepo(ddbClient)

	mediaBaseURL := os.Getenv("MEDIA_PRESIGN_BASE_URL")
	if mediaBaseURL == "" {
		mediaBaseURL = "https://journey-media.seonology.local"
	}
	mediaRepo := media.NewDDBRepo(ddbClient)
	mediaSvc := media.NewService(mediaRepo, media.NewStubPresigner(mediaBaseURL))

	deps := server.Deps{
		Trip:          trip.NewService(tripRepo),
		Day:           day.NewService(dayRepo),
		Schedule:      schedule.NewService(scheduleRepo),
		Meal:          meal.NewService(mealRepo),
		Accommodation: accommodation.NewService(accomRepo),
		ScheduleRepo:  scheduleRepo,
		MealRepo:      mealRepo,
		AccommRepo:    accomRepo,

		Expense:         record.NewExpenseService(expRepo),
		Note:            record.NewNoteService(noteRepo),
		Checklist:       record.NewChecklistService(chkRepo),
		Reservation:     record.NewReservationService(rsvRepo),
		ExpenseRepo:     expRepo,
		NoteRepo:        noteRepo,
		ChecklistRepo:   chkRepo,
		ReservationRepo: rsvRepo,

		Companion:     social.NewCompanionService(compRepo),
		Tag:           social.NewTagService(tagRepo),
		Template:      social.NewTemplateService(tplRepo),
		Favorite:      social.NewFavoriteService(favRepo),
		Share:         social.NewShareService(shareRepo),
		CompanionRepo: compRepo,
		TagRepo:       tagRepo,

		Media:     mediaSvc,
		MediaRepo: mediaRepo,
	}
	journey := server.NewJourneyServer(deps)

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			observability.LoggingUnaryInterceptor(logger),
		),
	)
	journeyv1.RegisterJourneyServiceServer(srv, journey)
	healthgrpc.RegisterHealthServer(srv, health.NewServer())
	reflection.Register(srv)

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("seonology-journey-back gRPC listening", "addr", addr)
		if err := srv.Serve(lis); err != nil {
			logger.Error("gRPC serve failed", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down")
	srv.GracefulStop()
	_ = obsServer.Close()
}
