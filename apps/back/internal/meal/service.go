// Package meal - Meal 도메인 (Day x MealType 조합으로 1).
package meal

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ErrNotFound - 미존재.
var ErrNotFound = errors.New("meal: not found")

func key(dayID string, mt journeyv1.MealType) string {
	return fmt.Sprintf("%s#%d", dayID, mt)
}

// Repository - Meal 저장소.
type Repository interface {
	Upsert(ctx context.Context, m *journeyv1.Meal) error
	ListByDay(ctx context.Context, dayID string) ([]*journeyv1.Meal, error)
	Delete(ctx context.Context, dayID string, mt journeyv1.MealType) error
	DeleteByDay(ctx context.Context, dayID string) error
}

// MemoryRepo - 메모리 Repository.
type MemoryRepo struct {
	mu    sync.RWMutex
	store map[string]*journeyv1.Meal
}

// NewMemoryRepo - 생성.
func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{store: make(map[string]*journeyv1.Meal)}
}

// Upsert - 생성 또는 갱신.
func (r *MemoryRepo) Upsert(_ context.Context, m *journeyv1.Meal) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[key(m.GetDayId(), m.GetMealType())] = clone(m)
	return nil
}

// ListByDay - day 의 식사 목록.
func (r *MemoryRepo) ListByDay(_ context.Context, dayID string) ([]*journeyv1.Meal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*journeyv1.Meal, 0)
	for _, m := range r.store {
		if m.GetDayId() == dayID {
			out = append(out, clone(m))
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].GetMealType() < out[j].GetMealType() })
	return out, nil
}

// Delete - 단일 삭제.
func (r *MemoryRepo) Delete(_ context.Context, dayID string, mt journeyv1.MealType) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	k := key(dayID, mt)
	if _, ok := r.store[k]; !ok {
		return ErrNotFound
	}
	delete(r.store, k)
	return nil
}

// DeleteByDay - cascade.
func (r *MemoryRepo) DeleteByDay(_ context.Context, dayID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for k, m := range r.store {
		if m.GetDayId() == dayID {
			delete(r.store, k)
		}
	}
	return nil
}

func clone(m *journeyv1.Meal) *journeyv1.Meal {
	if m == nil {
		return nil
	}
	cp := *m
	if m.GetAudit() != nil {
		ad := *m.Audit
		cp.Audit = &ad
	}
	if m.GetCost() != nil {
		c := *m.Cost
		cp.Cost = &c
	}
	if m.GetLocation() != nil {
		l := *m.Location
		cp.Location = &l
	}
	return &cp
}

// Service - Meal 로직.
type Service struct {
	repo Repository
	now  func() time.Time
}

// NewService - 생성.
func NewService(repo Repository) *Service {
	return &Service{repo: repo, now: time.Now}
}

// Upsert - 생성/갱신.
func (s *Service) Upsert(ctx context.Context, req *journeyv1.UpsertMealRequest) (*journeyv1.Meal, error) {
	if req.GetDayId() == "" || req.GetMealType() == journeyv1.MealType_MEAL_TYPE_UNSPECIFIED {
		return nil, errors.New("meal: day_id and meal_type required")
	}
	now := timestamppb.New(s.now().UTC())
	m := &journeyv1.Meal{
		DayId:          req.GetDayId(),
		MealType:       req.GetMealType(),
		Source:         req.GetSource(),
		RestaurantName: req.GetRestaurantName(),
		Menu:           req.GetMenu(),
		Cost:           req.GetCost(),
		Rating:         req.GetRating(),
		Review:         req.GetReview(),
		Location:       req.GetLocation(),
		Audit:          &journeyv1.AuditTimestamps{CreatedAt: now, UpdatedAt: now},
	}
	if err := s.repo.Upsert(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

// List - 일자별 목록.
func (s *Service) List(ctx context.Context, dayID string) ([]*journeyv1.Meal, error) {
	return s.repo.ListByDay(ctx, dayID)
}

// Delete - 삭제.
func (s *Service) Delete(ctx context.Context, dayID string, mt journeyv1.MealType) error {
	return s.repo.Delete(ctx, dayID, mt)
}
