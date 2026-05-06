// Package trip - Trip 도메인 (여행 전체 메타).
package trip

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

// ErrNotFound - Trip 미존재.
var ErrNotFound = errors.New("trip: not found")

// ErrForbidden - 소유자 아님.
var ErrForbidden = errors.New("trip: forbidden")

// Repository - Trip 저장소 인터페이스. 향후 DDB 어댑터로 교체.
type Repository interface {
	Create(ctx context.Context, t *journeyv1.Trip) error
	Get(ctx context.Context, id string) (*journeyv1.Trip, error)
	ListByOwner(ctx context.Context, ownerID string, status journeyv1.TripStatus) ([]*journeyv1.Trip, error)
	ListByOwnerPage(ctx context.Context, ownerID string, status journeyv1.TripStatus, cursor string, limit int32) ([]*journeyv1.Trip, string, error)
	Update(ctx context.Context, t *journeyv1.Trip) error
	Delete(ctx context.Context, id string) error
}

// MemoryRepo - 메모리 기반 Repository (개발/테스트용).
type MemoryRepo struct {
	mu    sync.RWMutex
	store map[string]*journeyv1.Trip
}

// NewMemoryRepo - MemoryRepo 생성.
func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{store: make(map[string]*journeyv1.Trip)}
}

// Create - 신규 Trip 저장.
func (r *MemoryRepo) Create(_ context.Context, t *journeyv1.Trip) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[t.GetId()] = clone(t)
	return nil
}

// Get - id 로 조회.
func (r *MemoryRepo) Get(_ context.Context, id string) (*journeyv1.Trip, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.store[id]
	if !ok {
		return nil, ErrNotFound
	}
	return clone(t), nil
}

// ListByOwner - Trip 목록 (최근 생성 순). ownerID 무시 — 전 유저 공유.
func (r *MemoryRepo) ListByOwner(_ context.Context, _ string, status journeyv1.TripStatus) ([]*journeyv1.Trip, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*journeyv1.Trip, 0, len(r.store))
	for _, t := range r.store {
		if status != journeyv1.TripStatus_TRIP_STATUS_UNSPECIFIED && t.GetStatus() != status {
			continue
		}
		out = append(out, clone(t))
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].GetAudit().GetCreatedAt().AsTime().After(out[j].GetAudit().GetCreatedAt().AsTime())
	})
	return out, nil
}

// ListByOwnerPage - MemoryRepo 는 in-process 용이라 cursor 를 사용하지 않고
// 단순 slice 기반 pagination 을 제공한다. DDB 구현체는 실제 cursor 를 사용.
func (r *MemoryRepo) ListByOwnerPage(ctx context.Context, ownerID string, status journeyv1.TripStatus, cursor string, limit int32) ([]*journeyv1.Trip, string, error) {
	all, err := r.ListByOwner(ctx, ownerID, status)
	if err != nil {
		return nil, "", err
	}
	offset := 0
	if cursor != "" {
		for i, t := range all {
			if t.GetId() == cursor {
				offset = i + 1
				break
			}
		}
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	end := offset + int(limit)
	if end > len(all) {
		end = len(all)
	}
	page := all[offset:end]
	next := ""
	if end < len(all) && len(page) > 0 {
		next = page[len(page)-1].GetId()
	}
	return page, next, nil
}

// Update - 기존 Trip 갱신.
func (r *MemoryRepo) Update(_ context.Context, t *journeyv1.Trip) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[t.GetId()]; !ok {
		return ErrNotFound
	}
	r.store[t.GetId()] = clone(t)
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

func clone(t *journeyv1.Trip) *journeyv1.Trip {
	if t == nil {
		return nil
	}
	return proto.Clone(t).(*journeyv1.Trip)
}

// Service - Trip 비즈니스 로직.
type Service struct {
	repo Repository
	now  func() time.Time
}

// NewService - Service 생성.
func NewService(repo Repository) *Service {
	return &Service{repo: repo, now: time.Now}
}

// Create - 신규 Trip.
func (s *Service) Create(ctx context.Context, ownerID string, req *journeyv1.CreateTripRequest) (*journeyv1.Trip, error) {
	now := timestamppb.New(s.now().UTC())
	t := &journeyv1.Trip{
		Id:          uuid.NewString(),
		OwnerId:     ownerID,
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		StartDate:   req.GetStartDate(),
		EndDate:     req.GetEndDate(),
		Status:      journeyv1.TripStatus_TRIP_STATUS_PLANNING,
		Destination: req.GetDestination(),
		CountryCode: req.GetCountryCode(),
		TotalBudget: req.GetTotalBudget(),
		Audit:       &journeyv1.AuditTimestamps{CreatedAt: now, UpdatedAt: now},
	}
	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

// Get - 조회. 소유자 무관 — 전 유저 공유.
func (s *Service) Get(ctx context.Context, _, id string) (*journeyv1.Trip, error) {
	return s.repo.Get(ctx, id)
}

// List - owner 의 trip 목록.
func (s *Service) List(ctx context.Context, ownerID string, status journeyv1.TripStatus) ([]*journeyv1.Trip, error) {
	return s.repo.ListByOwner(ctx, ownerID, status)
}

// ListPage - owner 의 trip 목록 (cursor 기반 페이지네이션).
// limit: 1~100 (0/음수/초과치는 20 으로 치환).
func (s *Service) ListPage(ctx context.Context, ownerID string, status journeyv1.TripStatus, cursor string, limit int32) ([]*journeyv1.Trip, string, error) {
	return s.repo.ListByOwnerPage(ctx, ownerID, status, cursor, limit)
}

// Update - 갱신. 소유자 무관 — 전 유저 공유.
func (s *Service) Update(ctx context.Context, _ string, req *journeyv1.UpdateTripRequest) (*journeyv1.Trip, error) {
	t, err := s.repo.Get(ctx, req.GetTripId())
	if err != nil {
		return nil, err
	}
	if req.GetTitle() != "" {
		t.Title = req.GetTitle()
	}
	if req.GetDescription() != "" {
		t.Description = req.GetDescription()
	}
	if req.GetStartDate() != "" {
		t.StartDate = req.GetStartDate()
	}
	if req.GetEndDate() != "" {
		t.EndDate = req.GetEndDate()
	}
	if req.GetStatus() != journeyv1.TripStatus_TRIP_STATUS_UNSPECIFIED {
		t.Status = req.GetStatus()
	}
	if req.GetCoverImageUrl() != "" {
		t.CoverImageUrl = req.GetCoverImageUrl()
	}
	if req.GetTotalBudget() != nil {
		t.TotalBudget = req.GetTotalBudget()
	}
	if req.GetDestination() != "" {
		t.Destination = req.GetDestination()
	}
	if req.GetCountryCode() != "" {
		t.CountryCode = req.GetCountryCode()
	}
	if t.Audit == nil {
		t.Audit = &journeyv1.AuditTimestamps{}
	}
	t.Audit.UpdatedAt = timestamppb.New(s.now().UTC())
	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

// Delete - 삭제. 소유자 무관 — 전 유저 공유.
func (s *Service) Delete(ctx context.Context, _, id string) error {
	return s.repo.Delete(ctx, id)
}
