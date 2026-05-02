package server

import (
	"context"
	"time"

	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// === Statistics ===

// GetTripStatistics implements JourneyService.
//
// MVP: 메모리 repo 들을 직접 집계.
func (s *JourneyServer) GetTripStatistics(ctx context.Context, req *journeyv1.GetTripStatisticsRequest) (*journeyv1.GetTripStatisticsResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.requireTripOwner(ctx, owner, req.GetTripId()); err != nil {
		return nil, err
	}
	days, _ := s.d.Day.List(ctx, req.GetTripId())
	stats := &journeyv1.TripStatistics{TripId: req.GetTripId(), TotalDays: int32(len(days))}
	regions := map[string]struct{}{}
	for _, d := range days {
		if d.GetRegion() != "" {
			regions[d.GetRegion()] = struct{}{}
		}
		if scs, err := s.d.Schedule.List(ctx, d.GetId()); err == nil {
			stats.TotalSchedules += int32(len(scs))
		}
		if ms, err := s.d.Meal.List(ctx, d.GetId()); err == nil {
			stats.TotalMeals += int32(len(ms))
		}
	}
	stats.VisitedRegions = int32(len(regions))
	if s.d.Expense != nil {
		if sum, err := s.d.Expense.Summary(ctx, req.GetTripId()); err == nil {
			stats.TotalExpense = sum.GetGrandTotal()
		}
	}
	if s.d.MediaRepo != nil {
		if n, err := s.d.MediaRepo.CountByTrip(ctx, req.GetTripId()); err == nil {
			stats.TotalPhotos = int32(n)
		}
	}
	return &journeyv1.GetTripStatisticsResponse{Stats: stats}, nil
}

// GetYearlyStatistics implements JourneyService - MVP stub.
func (s *JourneyServer) GetYearlyStatistics(ctx context.Context, req *journeyv1.GetYearlyStatisticsRequest) (*journeyv1.GetYearlyStatisticsResponse, error) {
	if _, err := ownerFromCtx(ctx); err != nil {
		return nil, err
	}
	return &journeyv1.GetYearlyStatisticsResponse{Stats: &journeyv1.YearlyStatistics{Year: req.GetYear()}}, nil
}

// === External (stubs) ===

// Geocode implements JourneyService.
func (s *JourneyServer) Geocode(_ context.Context, _ *journeyv1.GeocodeRequest) (*journeyv1.GeocodeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "geocode not configured")
}

// ReverseGeocode implements JourneyService.
func (s *JourneyServer) ReverseGeocode(_ context.Context, _ *journeyv1.ReverseGeocodeRequest) (*journeyv1.ReverseGeocodeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "reverse geocode not configured")
}

// GetExchangeRate implements JourneyService - 정적 stub (KRW->JPY=0.11 등).
func (s *JourneyServer) GetExchangeRate(_ context.Context, req *journeyv1.GetExchangeRateRequest) (*journeyv1.GetExchangeRateResponse, error) {
	rate := 1.0
	switch req.GetBase() + "->" + req.GetTarget() {
	case "KRW->JPY":
		rate = 0.11
	case "JPY->KRW":
		rate = 9.0
	case "KRW->USD":
		rate = 0.00072
	case "USD->KRW":
		rate = 1380.0
	}
	return &journeyv1.GetExchangeRateResponse{
		Base:      req.GetBase(),
		Target:    req.GetTarget(),
		Rate:      rate,
		FetchedAt: timestamppb.New(time.Now().UTC()),
	}, nil
}

// GetWeatherForecast implements JourneyService.
func (s *JourneyServer) GetWeatherForecast(_ context.Context, _ *journeyv1.GetWeatherForecastRequest) (*journeyv1.GetWeatherForecastResponse, error) {
	return nil, status.Error(codes.Unimplemented, "weather not configured")
}

// === Realtime (stub) ===

// PublishEvent implements JourneyService - hub 미연결 환경에서는 noop.
func (s *JourneyServer) PublishEvent(_ context.Context, _ *journeyv1.PublishEventRequest) (*journeyv1.PublishEventResponse, error) {
	return &journeyv1.PublishEventResponse{}, nil
}
