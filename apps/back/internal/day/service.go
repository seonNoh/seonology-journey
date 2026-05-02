// Package day - Day 도메인 (여행 일차).
package day

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ErrNotFound - Day 미존재.
var ErrNotFound = errors.New("day: not found")

// Repository - Day 저장소 인터페이스.
type Repository interface {
	Create(ctx context.Context, d *journeyv1.Day) error
	Get(ctx context.Context, id string) (*journeyv1.Day, error)
	ListByTrip(ctx context.Context, tripID string) ([]*journeyv1.Day, error)
	Update(ctx context.Context, d *journeyv1.Day) error
	DeleteByTrip(ctx context.Context, tripID string) error
}

// MemoryRepo - 메모리 기반 Day Repository.
type MemoryRepo struct {
	mu    sync.RWMutex
	store map[string]*journeyv1.Day
}

// NewMemoryRepo - MemoryRepo 생성.
func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{store: make(map[string]*journeyv1.Day)}
}

// Create - 신규 Day.
func (r *MemoryRepo) Create(_ context.Context, d *journeyv1.Day) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[d.GetId()] = clone(d)
	return nil
}

// Get - id 조회.
func (r *MemoryRepo) Get(_ context.Context, id string) (*journeyv1.Day, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.store[id]
	if !ok {
		return nil, ErrNotFound
	}
	return clone(d), nil
}

// ListByTrip - trip 의 day 목록 (day_number 오름차순).
func (r *MemoryRepo) ListByTrip(_ context.Context, tripID string) ([]*journeyv1.Day, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*journeyv1.Day, 0)
	for _, d := range r.store {
		if d.GetTripId() == tripID {
			out = append(out, clone(d))
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].GetDayNumber() < out[j].GetDayNumber() })
	return out, nil
}

// Update - 갱신.
func (r *MemoryRepo) Update(_ context.Context, d *journeyv1.Day) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[d.GetId()]; !ok {
		return ErrNotFound
	}
	r.store[d.GetId()] = clone(d)
	return nil
}

// DeleteByTrip - trip 삭제 시 cascade.
func (r *MemoryRepo) DeleteByTrip(_ context.Context, tripID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, d := range r.store {
		if d.GetTripId() == tripID {
			delete(r.store, id)
		}
	}
	return nil
}

func clone(d *journeyv1.Day) *journeyv1.Day {
	if d == nil {
		return nil
	}
	return proto.Clone(d).(*journeyv1.Day)
}

// Service - Day 비즈니스 로직.
type Service struct {
	repo Repository
	now  func() time.Time
}

// NewService - Service 생성.
func NewService(repo Repository) *Service {
	return &Service{repo: repo, now: time.Now}
}

// dateLayout - YYYY-MM-DD.
const dateLayout = "2006-01-02"

var weekdayKR = map[time.Weekday]string{
	time.Sunday: "일", time.Monday: "월", time.Tuesday: "화",
	time.Wednesday: "수", time.Thursday: "목", time.Friday: "금", time.Saturday: "토",
}

// GenerateForTrip - Trip 의 날짜 범위로 Day 들을 자동 생성. 既存があれば全削除して再生成.
func (s *Service) GenerateForTrip(ctx context.Context, tripID, startDate, endDate string) ([]*journeyv1.Day, error) {
	start, err := time.Parse(dateLayout, startDate)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse(dateLayout, endDate)
	if err != nil {
		return nil, err
	}
	if end.Before(start) {
		return nil, errors.New("day: end_date before start_date")
	}
	if err := s.repo.DeleteByTrip(ctx, tripID); err != nil {
		return nil, err
	}
	now := timestamppb.New(s.now().UTC())
	days := make([]*journeyv1.Day, 0)
	num := int32(1)
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		day := &journeyv1.Day{
			Id:        uuid.NewString(),
			TripId:    tripID,
			DayNumber: num,
			Date:      d.Format(dateLayout),
			DayOfWeek: weekdayKR[d.Weekday()],
			Audit:     &journeyv1.AuditTimestamps{CreatedAt: now, UpdatedAt: now},
		}
		if err := s.repo.Create(ctx, day); err != nil {
			return nil, err
		}
		days = append(days, day)
		num++
	}
	return days, nil
}

// List - trip 의 day 목록.
func (s *Service) List(ctx context.Context, tripID string) ([]*journeyv1.Day, error) {
	return s.repo.ListByTrip(ctx, tripID)
}

// Get - id 조회.
func (s *Service) Get(ctx context.Context, id string) (*journeyv1.Day, error) {
	return s.repo.Get(ctx, id)
}

// Update - region/weather/daily_summary 만 갱신.
func (s *Service) Update(ctx context.Context, req *journeyv1.UpdateDayRequest) (*journeyv1.Day, error) {
	d, err := s.repo.Get(ctx, req.GetDayId())
	if err != nil {
		return nil, err
	}
	if req.GetRegion() != "" {
		d.Region = req.GetRegion()
	}
	if req.GetWeather() != "" {
		d.Weather = req.GetWeather()
	}
	if req.GetDailySummary() != "" {
		d.DailySummary = req.GetDailySummary()
	}
	if d.Audit == nil {
		d.Audit = &journeyv1.AuditTimestamps{}
	}
	d.Audit.UpdatedAt = timestamppb.New(s.now().UTC())
	if err := s.repo.Update(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

// DeleteByTrip - cascade.
func (s *Service) DeleteByTrip(ctx context.Context, tripID string) error {
	return s.repo.DeleteByTrip(ctx, tripID)
}
