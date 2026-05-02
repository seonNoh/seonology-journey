// Package accommodation - Accommodation 도메인 (Day 당 1).
package accommodation

import (
	"context"
	"errors"
	"sync"
	"time"

	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ErrNotFound - 미존재.
var ErrNotFound = errors.New("accommodation: not found")

// Repository - Accommodation 저장소.
type Repository interface {
	Upsert(ctx context.Context, a *journeyv1.Accommodation) error
	Get(ctx context.Context, dayID string) (*journeyv1.Accommodation, error)
	Delete(ctx context.Context, dayID string) error
}

// MemoryRepo - 메모리 Repository.
type MemoryRepo struct {
	mu    sync.RWMutex
	store map[string]*journeyv1.Accommodation
}

// NewMemoryRepo - 생성.
func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{store: make(map[string]*journeyv1.Accommodation)}
}

// Upsert - 생성/갱신.
func (r *MemoryRepo) Upsert(_ context.Context, a *journeyv1.Accommodation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[a.GetDayId()] = clone(a)
	return nil
}

// Get - 조회.
func (r *MemoryRepo) Get(_ context.Context, dayID string) (*journeyv1.Accommodation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.store[dayID]
	if !ok {
		return nil, ErrNotFound
	}
	return clone(a), nil
}

// Delete - 삭제.
func (r *MemoryRepo) Delete(_ context.Context, dayID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[dayID]; !ok {
		return ErrNotFound
	}
	delete(r.store, dayID)
	return nil
}

func clone(a *journeyv1.Accommodation) *journeyv1.Accommodation {
	if a == nil {
		return nil
	}
	return proto.Clone(a).(*journeyv1.Accommodation)
}

// Service - Accommodation 로직.
type Service struct {
	repo Repository
	now  func() time.Time
}

// NewService - 생성.
func NewService(repo Repository) *Service {
	return &Service{repo: repo, now: time.Now}
}

// Upsert - 등록/갱신.
func (s *Service) Upsert(ctx context.Context, req *journeyv1.UpsertAccommodationRequest) (*journeyv1.Accommodation, error) {
	if req.GetDayId() == "" {
		return nil, errors.New("accommodation: day_id required")
	}
	now := timestamppb.New(s.now().UTC())
	a := &journeyv1.Accommodation{
		DayId:         req.GetDayId(),
		Name:          req.GetName(),
		CheckInTime:   req.GetCheckInTime(),
		CheckOutTime:  req.GetCheckOutTime(),
		Cost:          req.GetCost(),
		Amenities:     req.GetAmenities(),
		Address:       req.GetAddress(),
		Location:      req.GetLocation(),
		ReservationId: req.GetReservationId(),
		Audit:         &journeyv1.AuditTimestamps{CreatedAt: now, UpdatedAt: now},
	}
	if err := s.repo.Upsert(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

// Get - 조회.
func (s *Service) Get(ctx context.Context, dayID string) (*journeyv1.Accommodation, error) {
	return s.repo.Get(ctx, dayID)
}

// Delete - 삭제.
func (s *Service) Delete(ctx context.Context, dayID string) error {
	return s.repo.Delete(ctx, dayID)
}
