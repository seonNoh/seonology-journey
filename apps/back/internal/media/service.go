// Package media - Media (사진/영상) 메타 + 업로드 URL 관리.
//
// MVP: presigned URL 은 내부 stub 으로 반환. 추후 S3 어댑터로 대체.
package media

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ErrNotFound - 미존재.
var ErrNotFound = errors.New("media: not found")

// Repository - 메타 저장소.
type Repository interface {
	Save(ctx context.Context, m *journeyv1.Media) error
	Get(ctx context.Context, id string) (*journeyv1.Media, error)
	ListByTrip(ctx context.Context, tripID, dayID string) ([]*journeyv1.Media, error)
	Delete(ctx context.Context, id string) error
	DeleteByTrip(ctx context.Context, tripID string) error
	CountByTrip(ctx context.Context, tripID string) (int, error)
}

// MemoryRepo - 메모리 메타 Repository.
type MemoryRepo struct {
	mu    sync.RWMutex
	store map[string]*journeyv1.Media
}

// NewMemoryRepo - 생성.
func NewMemoryRepo() *MemoryRepo { return &MemoryRepo{store: make(map[string]*journeyv1.Media)} }

// Save - 저장.
func (r *MemoryRepo) Save(_ context.Context, m *journeyv1.Media) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[m.GetId()] = m
	return nil
}

// Get - 조회.
func (r *MemoryRepo) Get(_ context.Context, id string) (*journeyv1.Media, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.store[id]
	if !ok {
		return nil, ErrNotFound
	}
	return m, nil
}

// ListByTrip - 목록 (dayID 비어있으면 trip 전체).
func (r *MemoryRepo) ListByTrip(_ context.Context, tripID, dayID string) ([]*journeyv1.Media, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*journeyv1.Media, 0)
	for _, m := range r.store {
		if m.GetTripId() != tripID {
			continue
		}
		if dayID != "" && m.GetDayId() != dayID {
			continue
		}
		out = append(out, m)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].GetTakenAt().AsTime().Before(out[j].GetTakenAt().AsTime())
	})
	return out, nil
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

// DeleteByTrip - cascade.
func (r *MemoryRepo) DeleteByTrip(_ context.Context, tripID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for k, m := range r.store {
		if m.GetTripId() == tripID {
			delete(r.store, k)
		}
	}
	return nil
}

// CountByTrip - 개수.
func (r *MemoryRepo) CountByTrip(_ context.Context, tripID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	n := 0
	for _, m := range r.store {
		if m.GetTripId() == tripID {
			n++
		}
	}
	return n, nil
}

// Presigner - 외부 업로드 URL 생성.
type Presigner interface {
	UploadURL(ctx context.Context, key, mime string, size int64) (url string, expiresAt time.Time, err error)
	DownloadURL(ctx context.Context, key string) (url string, expiresAt time.Time, err error)
}

// StubPresigner - 개발용 stub.
type StubPresigner struct {
	BaseURL string // ex: https://media.seonology-journey.local
	TTL     time.Duration
	now     func() time.Time
}

// NewStubPresigner - 생성.
func NewStubPresigner(base string) *StubPresigner {
	return &StubPresigner{BaseURL: base, TTL: 15 * time.Minute, now: time.Now}
}

// UploadURL - stub.
func (p *StubPresigner) UploadURL(_ context.Context, key, _ string, _ int64) (string, time.Time, error) {
	exp := p.now().Add(p.TTL)
	return fmt.Sprintf("%s/upload/%s?expires=%d", p.BaseURL, key, exp.Unix()), exp, nil
}

// DownloadURL - stub.
func (p *StubPresigner) DownloadURL(_ context.Context, key string) (string, time.Time, error) {
	exp := p.now().Add(p.TTL)
	return fmt.Sprintf("%s/get/%s?expires=%d", p.BaseURL, key, exp.Unix()), exp, nil
}

// Service - 미디어 로직.
type Service struct {
	repo  Repository
	pres  Presigner
	now   func() time.Time
	keyer func(tripID, mediaID, filename string) string
}

// NewService - 생성.
func NewService(repo Repository, pres Presigner) *Service {
	return &Service{
		repo:  repo,
		pres:  pres,
		now:   time.Now,
		keyer: func(tripID, mediaID, filename string) string { return fmt.Sprintf("trips/%s/%s/%s", tripID, mediaID, filename) },
	}
}

// GetUploadURL - presigned PUT URL 발급. media 메타는 confirm 시 저장.
func (s *Service) GetUploadURL(ctx context.Context, req *journeyv1.GetUploadUrlRequest) (*journeyv1.GetUploadUrlResponse, error) {
	mediaID := uuid.NewString()
	key := s.keyer(req.GetTripId(), mediaID, req.GetFilename())
	url, exp, err := s.pres.UploadURL(ctx, key, req.GetMimeType(), req.GetSize())
	if err != nil {
		return nil, err
	}
	return &journeyv1.GetUploadUrlResponse{
		UploadUrl: url,
		S3Key:     key,
		ExpiresAt: timestamppb.New(exp),
		MediaId:   mediaID,
	}, nil
}

// ConfirmUpload - 업로드 완료 후 메타 저장.
func (s *Service) ConfirmUpload(ctx context.Context, req *journeyv1.ConfirmUploadRequest) (*journeyv1.Media, error) {
	now := timestamppb.New(s.now().UTC())
	m := &journeyv1.Media{
		Id:         req.GetMediaId(),
		TripId:     req.GetTripId(),
		DayId:      req.GetDayId(),
		ScheduleId: req.GetScheduleId(),
		S3Key:      req.GetS3Key(),
		Caption:    req.GetCaption(),
		TakenAt:    now,
		Audit:      &journeyv1.AuditTimestamps{CreatedAt: now, UpdatedAt: now},
	}
	return m, s.repo.Save(ctx, m)
}

// List - trip / day 목록.
func (s *Service) List(ctx context.Context, tripID, dayID string) ([]*journeyv1.Media, error) {
	return s.repo.ListByTrip(ctx, tripID, dayID)
}

// Get - id 조회.
func (s *Service) Get(ctx context.Context, id string) (*journeyv1.Media, error) {
	return s.repo.Get(ctx, id)
}

// Delete - 삭제.
func (s *Service) Delete(ctx context.Context, id string) error { return s.repo.Delete(ctx, id) }

// URL - presigned GET.
func (s *Service) URL(ctx context.Context, id string, thumbnail bool) (*journeyv1.GetMediaUrlResponse, error) {
	m, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	key := m.GetS3Key()
	if thumbnail && m.GetThumbnailS3Key() != "" {
		key = m.GetThumbnailS3Key()
	}
	url, exp, err := s.pres.DownloadURL(ctx, key)
	if err != nil {
		return nil, err
	}
	return &journeyv1.GetMediaUrlResponse{Url: url, ExpiresAt: timestamppb.New(exp)}, nil
}
