// Package main is the entrypoint for the seonology-journey-back gRPC service.
package main

import (
	"log"
	"net"
	"os"

	"github.com/seonNoh/seonology-journey/apps/back/internal/accommodation"
	"github.com/seonNoh/seonology-journey/apps/back/internal/day"
	"github.com/seonNoh/seonology-journey/apps/back/internal/meal"
	"github.com/seonNoh/seonology-journey/apps/back/internal/schedule"
	"github.com/seonNoh/seonology-journey/apps/back/internal/server"
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

	deps := server.Deps{
		Trip:          trip.NewService(trip.NewMemoryRepo()),
		Day:           day.NewService(day.NewMemoryRepo()),
		Schedule:      schedule.NewService(scheduleRepo),
		Meal:          meal.NewService(mealRepo),
		Accommodation: accommodation.NewService(accomRepo),
		ScheduleRepo:  scheduleRepo,
		MealRepo:      mealRepo,
		AccommRepo:    accomRepo,
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
