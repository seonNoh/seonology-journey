// Package schedule - Schedule 도메인 (Day 내 일정).
package schedule

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ErrNotFound - 미존재.
var ErrNotFound = errors.New("schedule: not found")

// Repository - Schedule 저장소 인터페이스.
type Repository interface {
	Create(ctx context.Context, s *journeyv1.Schedule) error
	Get(ctx context.Context, id string) (*journeyv1.Schedule, error)
	ListByDay(ctx context.Context, dayID string) ([]*journeyv1.Schedule, error)
	Update(ctx context.Context, s *journeyv1.Schedule) error
	Delete(ctx context.Context, id string) error
}

// MemoryRepo - 메모리 Repository.
type MemoryRepo struct {
	mu    sync.RWMutex
	store map[string]*journeyv1.Schedule
}

// NewMemoryRepo - 생성.
func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{store: make(map[string]*journeyv1.Schedule)}
}

// Create - 생성.
func (r *MemoryRepo) Create(_ context.Context, s *journeyv1.Schedule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[s.GetId()] = clone(s)
	return nil
}

// Get - 조회.
func (r *MemoryRepo) Get(_ context.Context, id string) (*journeyv1.Schedule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.store[id]
	if !ok {
		return nil, ErrNotFound
	}
	return clone(s), nil
}

// ListByDay - day 의 schedule 목록 (order 오름차순).
func (r *MemoryRepo) ListByDay(_ context.Context, dayID string) ([]*journeyv1.Schedule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*journeyv1.Schedule, 0)
	for _, s := range r.store {
		if s.GetDayId() == dayID {
			out = append(out, clone(s))
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].GetOrder() < out[j].GetOrder() })
	return out, nil
}

// Update - 갱신.
func (r *MemoryRepo) Update(_ context.Context, s *journeyv1.Schedule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[s.GetId()]; !ok {
		return ErrNotFound
	}
	r.store[s.GetId()] = clone(s)
	return nil
}

// Delete - 삭제.
func (r *MemoryRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[id]; !ok {
		return ErrNotFound
	}
	delete(r.store, id)
	return nil
}

func clone(s *journeyv1.Schedule) *journeyv1.Schedule {
	if s == nil {
		return nil
	}
	cp := *s
	if s.GetAudit() != nil {
		ad := *s.Audit
		cp.Audit = &ad
	}
	if s.GetCost() != nil {
		c := *s.Cost
		cp.Cost = &c
	}
	if s.GetLocation() != nil {
		l := *s.Location
		cp.Location = &l
	}
	return &cp
}

// Service - Schedule 로직.
type Service struct {
	repo Repository
	now  func() time.Time
}

// NewService - 생성.
func NewService(repo Repository) *Service {
	return &Service{repo: repo, now: time.Now}
}

// Create - 신규. order 는 마지막+1.
func (s *Service) Create(ctx context.Context, req *journeyv1.CreateScheduleRequest) (*journeyv1.Schedule, error) {
	if req.GetDayId() == "" {
		return nil, errors.New("schedule: day_id required")
	}
	existing, err := s.repo.ListByDay(ctx, req.GetDayId())
	if err != nil {
		return nil, err
	}
	now := timestamppb.New(s.now().UTC())
	sch := &journeyv1.Schedule{
		Id:              uuid.NewString(),
		DayId:           req.GetDayId(),
		Order:           int32(len(existing)) + 1,
		StartTime:       req.GetStartTime(),
		EndTime:         req.GetEndTime(),
		Title:           req.GetTitle(),
		Region:          req.GetRegion(),
		Category:        req.GetCategory(),
		Transport:       req.GetTransport(),
		TransportDetail: req.GetTransportDetail(),
		Cost:            req.GetCost(),
		PlaceName:       req.GetPlaceName(),
		Location:        req.GetLocation(),
		Notes:           req.GetNotes(),
		Audit:           &journeyv1.AuditTimestamps{CreatedAt: now, UpdatedAt: now},
	}
	if err := s.repo.Create(ctx, sch); err != nil {
		return nil, err
	}
	return sch, nil
}

// List - day 의 schedule 목록.
func (s *Service) List(ctx context.Context, dayID string) ([]*journeyv1.Schedule, error) {
	return s.repo.ListByDay(ctx, dayID)
}

// Update - 갱신.
func (s *Service) Update(ctx context.Context, req *journeyv1.UpdateScheduleRequest) (*journeyv1.Schedule, error) {
	sch, err := s.repo.Get(ctx, req.GetScheduleId())
	if err != nil {
		return nil, err
	}
	if req.GetStartTime() != "" {
		sch.StartTime = req.GetStartTime()
	}
	if req.GetEndTime() != "" {
		sch.EndTime = req.GetEndTime()
	}
	if req.GetTitle() != "" {
		sch.Title = req.GetTitle()
	}
	if req.GetRegion() != "" {
		sch.Region = req.GetRegion()
	}
	if req.GetCategory() != journeyv1.ScheduleCategory_SCHEDULE_CATEGORY_UNSPECIFIED {
		sch.Category = req.GetCategory()
	}
	if req.GetTransport() != journeyv1.TransportType_TRANSPORT_TYPE_UNSPECIFIED {
		sch.Transport = req.GetTransport()
	}
	if req.GetTransportDetail() != "" {
		sch.TransportDetail = req.GetTransportDetail()
	}
	if req.GetCost() != nil {
		sch.Cost = req.GetCost()
	}
	if req.GetPlaceName() != "" {
		sch.PlaceName = req.GetPlaceName()
	}
	if req.GetLocation() != nil {
		sch.Location = req.GetLocation()
	}
	if req.GetNotes() != "" {
		sch.Notes = req.GetNotes()
	}
	sch.IsCompleted = req.GetIsCompleted()
	if sch.Audit == nil {
		sch.Audit = &journeyv1.AuditTimestamps{}
	}
	sch.Audit.UpdatedAt = timestamppb.New(s.now().UTC())
	if err := s.repo.Update(ctx, sch); err != nil {
		return nil, err
	}
	return sch, nil
}

// Delete - 삭제.
func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// DeleteByDay - day 삭제 시 cascade.
func (r *MemoryRepo) DeleteByDay(_ context.Context, dayID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, s := range r.store {
		if s.GetDayId() == dayID {
			delete(r.store, id)
		}
	}
	return nil
}

// Reorder - 일괄 순서 변경.
func (s *Service) Reorder(ctx context.Context, dayID string, orderedIDs []string) ([]*journeyv1.Schedule, error) {
	now := timestamppb.New(s.now().UTC())
	for i, id := range orderedIDs {
		sch, err := s.repo.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		if sch.GetDayId() != dayID {
			return nil, errors.New("schedule: id does not belong to day_id")
		}
		sch.Order = int32(i) + 1
		if sch.Audit == nil {
			sch.Audit = &journeyv1.AuditTimestamps{}
		}
		sch.Audit.UpdatedAt = now
		if err := s.repo.Update(ctx, sch); err != nil {
			return nil, err
		}
	}
	return s.repo.ListByDay(ctx, dayID)
}
