// Package main is the entrypoint for the seonology-journey-back gRPC service.
package main

import (
	"log"
	"net"
	"os"

	"github.com/seonNoh/seonology-journey/apps/back/internal/accommodation"
	"github.com/seonNoh/seonology-journey/apps/back/internal/day"
	"github.com/seonNoh/seonology-journey/apps/back/internal/meal"
	"github.com/seonNoh/seonology-journey/apps/back/internal/media"
	"github.com/seonNoh/seonology-journey/apps/back/internal/record"
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
	addr := os.Getenv("GRPC_LISTEN_ADDR")
	if addr == "" {
		addr = ":9090"
	}
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	scheduleRepo := schedule.NewMemoryRepo()
	mealRepo := meal.NewMemoryRepo()
	accomRepo := accommodation.NewMemoryRepo()

	expRepo := record.NewExpenseRepo()
	noteRepo := record.NewNoteRepo()
	chkRepo := record.NewChecklistRepo()
	rsvRepo := record.NewReservationRepo()

	compRepo := social.NewCompanionRepo()
	tagRepo := social.NewTagRepo()
	tplRepo := social.NewTemplateRepo()
	favRepo := social.NewFavoriteRepo()
	shareRepo := social.NewShareRepo()

	mediaBaseURL := os.Getenv("MEDIA_PRESIGN_BASE_URL")
	if mediaBaseURL == "" {
		mediaBaseURL = "https://journey-media.seonology.local"
	}
	mediaRepo := media.NewMemoryRepo()
	mediaSvc := media.NewService(mediaRepo, media.NewStubPresigner(mediaBaseURL))

	deps := server.Deps{
		Trip:          trip.NewService(trip.NewMemoryRepo()),
		Day:           day.NewService(day.NewMemoryRepo()),
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

	srv := grpc.NewServer()
	journeyv1.RegisterJourneyServiceServer(srv, journey)
	healthgrpc.RegisterHealthServer(srv, health.NewServer())
	reflection.Register(srv)

	log.Printf("seonology-journey-back gRPC listening on %s", addr)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
