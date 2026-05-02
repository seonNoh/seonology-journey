// Package server - JourneyService gRPC 어댑터.
package server

import (
	"context"
	"errors"

	"github.com/seonNoh/seonology-journey/apps/back/internal/accommodation"
	"github.com/seonNoh/seonology-journey/apps/back/internal/day"
	"github.com/seonNoh/seonology-journey/apps/back/internal/meal"
	"github.com/seonNoh/seonology-journey/apps/back/internal/schedule"
	"github.com/seonNoh/seonology-journey/apps/back/internal/trip"
	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Deps - JourneyServer 의존성.
type Deps struct {
	Trip          *trip.Service
	Day           *day.Service
	Schedule      *schedule.Service
	Meal          *meal.Service
	Accommodation *accommodation.Service
	ScheduleRepo  *schedule.MemoryRepo
	MealRepo      *meal.MemoryRepo
	AccommRepo    *accommodation.MemoryRepo
}

// JourneyServer - JourneyService 구현.
type JourneyServer struct {
	journeyv1.UnimplementedJourneyServiceServer
	d Deps
}

// NewJourneyServer - 생성.
func NewJourneyServer(d Deps) *JourneyServer {
	return &JourneyServer{d: d}
}

// ownerFromCtx - x-user-id metadata.
func ownerFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}
	vals := md.Get("x-user-id")
	if len(vals) == 0 || vals[0] == "" {
		return "", status.Error(codes.Unauthenticated, "missing x-user-id")
	}
	return vals[0], nil
}

func mapErr(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, trip.ErrNotFound),
		errors.Is(err, day.ErrNotFound),
		errors.Is(err, schedule.ErrNotFound),
		errors.Is(err, meal.ErrNotFound),
		errors.Is(err, accommodation.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, trip.ErrForbidden):
		return status.Error(codes.PermissionDenied, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

// requireDayOwner - day -> trip -> owner 검증.
func (s *JourneyServer) requireDayOwner(ctx context.Context, owner, dayID string) (*journeyv1.Day, error) {
	d, err := s.d.Day.Get(ctx, dayID)
	if err != nil {
		return nil, mapErr(err)
	}
	if _, err := s.d.Trip.Get(ctx, owner, d.GetTripId()); err != nil {
		return nil, mapErr(err)
	}
	return d, nil
}

// requireScheduleOwner - schedule 의 owner 검증.
func (s *JourneyServer) requireScheduleOwner(ctx context.Context, owner, scheduleID string) (*journeyv1.Schedule, error) {
	sch, err := s.d.ScheduleRepo.Get(ctx, scheduleID)
	if err != nil {
		return nil, mapErr(err)
	}
	if _, err := s.requireDayOwner(ctx, owner, sch.GetDayId()); err != nil {
		return nil, err
	}
	return sch, nil
}

// === Trip ===

// CreateTrip implements JourneyService.
func (s *JourneyServer) CreateTrip(ctx context.Context, req *journeyv1.CreateTripRequest) (*journeyv1.CreateTripResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "title required")
	}
	t, err := s.d.Trip.Create(ctx, owner, req)
	if err != nil {
		return nil, mapErr(err)
	}
	if t.GetStartDate() != "" && t.GetEndDate() != "" {
		if _, err := s.d.Day.GenerateForTrip(ctx, t.GetId(), t.GetStartDate(), t.GetEndDate()); err != nil {
			return nil, mapErr(err)
		}
	}
	return &journeyv1.CreateTripResponse{Trip: t}, nil
}

// GetTrip implements JourneyService.
func (s *JourneyServer) GetTrip(ctx context.Context, req *journeyv1.GetTripRequest) (*journeyv1.GetTripResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	t, err := s.d.Trip.Get(ctx, owner, req.GetTripId())
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.GetTripResponse{Trip: t}, nil
}

// ListTrips implements JourneyService.
func (s *JourneyServer) ListTrips(ctx context.Context, req *journeyv1.ListTripsRequest) (*journeyv1.ListTripsResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	trips, err := s.d.Trip.List(ctx, owner, req.GetStatus())
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.ListTripsResponse{Trips: trips, Page: &journeyv1.PageInfo{HasMore: false}}, nil
}

// UpdateTrip implements JourneyService.
func (s *JourneyServer) UpdateTrip(ctx context.Context, req *journeyv1.UpdateTripRequest) (*journeyv1.UpdateTripResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	prev, err := s.d.Trip.Get(ctx, owner, req.GetTripId())
	if err != nil {
		return nil, mapErr(err)
	}
	t, err := s.d.Trip.Update(ctx, owner, req)
	if err != nil {
		return nil, mapErr(err)
	}
	if t.GetStartDate() != prev.GetStartDate() || t.GetEndDate() != prev.GetEndDate() {
		if _, err := s.d.Day.GenerateForTrip(ctx, t.GetId(), t.GetStartDate(), t.GetEndDate()); err != nil {
			return nil, mapErr(err)
		}
	}
	return &journeyv1.UpdateTripResponse{Trip: t}, nil
}

// DeleteTrip implements JourneyService.
func (s *JourneyServer) DeleteTrip(ctx context.Context, req *journeyv1.DeleteTripRequest) (*journeyv1.DeleteTripResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	days, _ := s.d.Day.List(ctx, req.GetTripId())
	if err := s.d.Trip.Delete(ctx, owner, req.GetTripId()); err != nil {
		return nil, mapErr(err)
	}
	for _, dd := range days {
		_ = s.d.ScheduleRepo.DeleteByDay(ctx, dd.GetId())
		_ = s.d.MealRepo.DeleteByDay(ctx, dd.GetId())
		_ = s.d.AccommRepo.Delete(ctx, dd.GetId())
	}
	_ = s.d.Day.DeleteByTrip(ctx, req.GetTripId())
	return &journeyv1.DeleteTripResponse{}, nil
}

// === Day ===

// ListDays implements JourneyService.
func (s *JourneyServer) ListDays(ctx context.Context, req *journeyv1.ListDaysRequest) (*journeyv1.ListDaysResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := s.d.Trip.Get(ctx, owner, req.GetTripId()); err != nil {
		return nil, mapErr(err)
	}
	days, err := s.d.Day.List(ctx, req.GetTripId())
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.ListDaysResponse{Days: days}, nil
}

// UpdateDay implements JourneyService.
func (s *JourneyServer) UpdateDay(ctx context.Context, req *journeyv1.UpdateDayRequest) (*journeyv1.UpdateDayResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := s.requireDayOwner(ctx, owner, req.GetDayId()); err != nil {
		return nil, err
	}
	d, err := s.d.Day.Update(ctx, req)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.UpdateDayResponse{Day: d}, nil
}

// === Schedule ===

// CreateSchedule implements JourneyService.
func (s *JourneyServer) CreateSchedule(ctx context.Context, req *journeyv1.CreateScheduleRequest) (*journeyv1.CreateScheduleResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := s.requireDayOwner(ctx, owner, req.GetDayId()); err != nil {
		return nil, err
	}
	sch, err := s.d.Schedule.Create(ctx, req)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.CreateScheduleResponse{Schedule: sch}, nil
}

// ListSchedules implements JourneyService.
func (s *JourneyServer) ListSchedules(ctx context.Context, req *journeyv1.ListSchedulesRequest) (*journeyv1.ListSchedulesResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := s.requireDayOwner(ctx, owner, req.GetDayId()); err != nil {
		return nil, err
	}
	out, err := s.d.Schedule.List(ctx, req.GetDayId())
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.ListSchedulesResponse{Schedules: out}, nil
}

// UpdateSchedule implements JourneyService.
func (s *JourneyServer) UpdateSchedule(ctx context.Context, req *journeyv1.UpdateScheduleRequest) (*journeyv1.UpdateScheduleResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := s.requireScheduleOwner(ctx, owner, req.GetScheduleId()); err != nil {
		return nil, err
	}
	sch, err := s.d.Schedule.Update(ctx, req)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.UpdateScheduleResponse{Schedule: sch}, nil
}

// DeleteSchedule implements JourneyService.
func (s *JourneyServer) DeleteSchedule(ctx context.Context, req *journeyv1.DeleteScheduleRequest) (*journeyv1.DeleteScheduleResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := s.requireScheduleOwner(ctx, owner, req.GetScheduleId()); err != nil {
		return nil, err
	}
	if err := s.d.Schedule.Delete(ctx, req.GetScheduleId()); err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.DeleteScheduleResponse{}, nil
}

// ReorderSchedules implements JourneyService.
func (s *JourneyServer) ReorderSchedules(ctx context.Context, req *journeyv1.ReorderSchedulesRequest) (*journeyv1.ReorderSchedulesResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := s.requireDayOwner(ctx, owner, req.GetDayId()); err != nil {
		return nil, err
	}
	out, err := s.d.Schedule.Reorder(ctx, req.GetDayId(), req.GetScheduleIdsInOrder())
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.ReorderSchedulesResponse{Schedules: out}, nil
}

// === Meal ===

// UpsertMeal implements JourneyService.
func (s *JourneyServer) UpsertMeal(ctx context.Context, req *journeyv1.UpsertMealRequest) (*journeyv1.UpsertMealResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := s.requireDayOwner(ctx, owner, req.GetDayId()); err != nil {
		return nil, err
	}
	m, err := s.d.Meal.Upsert(ctx, req)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.UpsertMealResponse{Meal: m}, nil
}

// ListMeals implements JourneyService.
func (s *JourneyServer) ListMeals(ctx context.Context, req *journeyv1.ListMealsRequest) (*journeyv1.ListMealsResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := s.requireDayOwner(ctx, owner, req.GetDayId()); err != nil {
		return nil, err
	}
	out, err := s.d.Meal.List(ctx, req.GetDayId())
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.ListMealsResponse{Meals: out}, nil
}

// DeleteMeal implements JourneyService.
func (s *JourneyServer) DeleteMeal(ctx context.Context, req *journeyv1.DeleteMealRequest) (*journeyv1.DeleteMealResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := s.requireDayOwner(ctx, owner, req.GetDayId()); err != nil {
		return nil, err
	}
	if err := s.d.Meal.Delete(ctx, req.GetDayId(), req.GetMealType()); err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.DeleteMealResponse{}, nil
}

// === Accommodation ===

// UpsertAccommodation implements JourneyService.
func (s *JourneyServer) UpsertAccommodation(ctx context.Context, req *journeyv1.UpsertAccommodationRequest) (*journeyv1.UpsertAccommodationResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := s.requireDayOwner(ctx, owner, req.GetDayId()); err != nil {
		return nil, err
	}
	a, err := s.d.Accommodation.Upsert(ctx, req)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.UpsertAccommodationResponse{Accommodation: a}, nil
}

// GetAccommodation implements JourneyService.
func (s *JourneyServer) GetAccommodation(ctx context.Context, req *journeyv1.GetAccommodationRequest) (*journeyv1.GetAccommodationResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := s.requireDayOwner(ctx, owner, req.GetDayId()); err != nil {
		return nil, err
	}
	a, err := s.d.Accommodation.Get(ctx, req.GetDayId())
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.GetAccommodationResponse{Accommodation: a}, nil
}

// DeleteAccommodation implements JourneyService.
func (s *JourneyServer) DeleteAccommodation(ctx context.Context, req *journeyv1.DeleteAccommodationRequest) (*journeyv1.DeleteAccommodationResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := s.requireDayOwner(ctx, owner, req.GetDayId()); err != nil {
		return nil, err
	}
	if err := s.d.Accommodation.Delete(ctx, req.GetDayId()); err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.DeleteAccommodationResponse{}, nil
}
