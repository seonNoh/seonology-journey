package server

import (
	"context"
	"sync"
	"time"

	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// === Statistics ===

// GetTripStatistics implements JourneyService.
//
// Uses a fan-out worker pool across days: each day contributes a
// schedule count, meal count, and region label. Fetches for distinct
// days run concurrently so overall latency is O(1) in the number of
// days plus a small expense/media aggregate pair, instead of O(N) as
// the previous implementation did.
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

	// Aggregations protected by a mutex; writes are small and coarse.
	var (
		mu            sync.Mutex
		scheduleTotal int32
		mealTotal     int32
		regions       = map[string]struct{}{}
	)

	// Cap concurrency so a trip with 30+ days doesn't spam DynamoDB.
	const maxParallel = 8
	sem := make(chan struct{}, maxParallel)
	var wg sync.WaitGroup

	for _, d := range days {
		d := d
		if d.GetRegion() != "" {
			mu.Lock()
			regions[d.GetRegion()] = struct{}{}
			mu.Unlock()
		}
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			scCount := int32(0)
			if scs, err := s.d.Schedule.List(ctx, d.GetId()); err == nil {
				scCount = int32(len(scs))
			}
			meCount := int32(0)
			if ms, err := s.d.Meal.List(ctx, d.GetId()); err == nil {
				meCount = int32(len(ms))
			}
			mu.Lock()
			scheduleTotal += scCount
			mealTotal += meCount
			mu.Unlock()
		}()
	}

	// Expense summary and media count run in parallel with the day fan-out.
	var (
		expenseWg      sync.WaitGroup
		expenseSummary *journeyv1.ExpenseSummary
		mediaCount     int
	)
	if s.d.Expense != nil {
		expenseWg.Add(1)
		go func() {
			defer expenseWg.Done()
			if sum, err := s.d.Expense.Summary(ctx, req.GetTripId()); err == nil {
				expenseSummary = sum
			}
		}()
	}
	if s.d.MediaRepo != nil {
		expenseWg.Add(1)
		go func() {
			defer expenseWg.Done()
			if n, err := s.d.MediaRepo.CountByTrip(ctx, req.GetTripId()); err == nil {
				mediaCount = n
			}
		}()
	}

	wg.Wait()
	expenseWg.Wait()

	stats.TotalSchedules = scheduleTotal
	stats.TotalMeals = mealTotal
	stats.VisitedRegions = int32(len(regions))
	stats.TotalPhotos = int32(mediaCount)
	if expenseSummary != nil {
		stats.TotalExpense = expenseSummary.GetGrandTotal()
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
