// Package social - Companion / Tag / Template / FavoritePlace / Share 도메인.
package social

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ErrNotFound - 미존재.
var ErrNotFound = errors.New("social: not found")

// === Companion ===

// CompanionRepo - 메모리.
type CompanionRepo struct {
	mu    sync.RWMutex
	store map[string]*journeyv1.Companion // key: tripID#memberID
}

// NewCompanionRepo - 생성.
func NewCompanionRepo() *CompanionRepo {
	return &CompanionRepo{store: make(map[string]*journeyv1.Companion)}
}

func compKey(tripID, memberID string) string { return tripID + "#" + memberID }

// Save - 저장.
func (r *CompanionRepo) Save(_ context.Context, c *journeyv1.Companion) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[compKey(c.GetTripId(), c.GetMemberId())] = c
	return nil
}

// Get - 조회.
func (r *CompanionRepo) Get(_ context.Context, tripID, memberID string) (*journeyv1.Companion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.store[compKey(tripID, memberID)]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

// ListByTrip - 목록.
func (r *CompanionRepo) ListByTrip(_ context.Context, tripID string) ([]*journeyv1.Companion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*journeyv1.Companion, 0)
	for _, c := range r.store {
		if c.GetTripId() == tripID {
			out = append(out, c)
		}
	}
	return out, nil
}

// Delete - 삭제.
func (r *CompanionRepo) Delete(_ context.Context, tripID, memberID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	k := compKey(tripID, memberID)
	if _, ok := r.store[k]; !ok {
		return ErrNotFound
	}
	delete(r.store, k)
	return nil
}

// DeleteByTrip - cascade.
func (r *CompanionRepo) DeleteByTrip(_ context.Context, tripID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for k, c := range r.store {
		if c.GetTripId() == tripID {
			delete(r.store, k)
		}
	}
	return nil
}

// CompanionService - 로직.
type CompanionService struct {
	repo CompanionRepository
	now  func() time.Time
}

// NewCompanionService - 생성.
func NewCompanionService(repo CompanionRepository) *CompanionService {
	return &CompanionService{repo: repo, now: time.Now}
}

// Add - 추가.
func (s *CompanionService) Add(ctx context.Context, req *journeyv1.AddCompanionRequest) (*journeyv1.Companion, error) {
	c := &journeyv1.Companion{
		TripId:    req.GetTripId(),
		MemberId:  req.GetMemberId(),
		Role:      req.GetRole(),
		InvitedAt: timestamppb.New(s.now().UTC()),
	}
	return c, s.repo.Save(ctx, c)
}

// List - 목록.
func (s *CompanionService) List(ctx context.Context, tripID string) ([]*journeyv1.Companion, error) {
	return s.repo.ListByTrip(ctx, tripID)
}

// UpdateRole - 권한 변경.
func (s *CompanionService) UpdateRole(ctx context.Context, req *journeyv1.UpdateCompanionRoleRequest) (*journeyv1.Companion, error) {
	c, err := s.repo.Get(ctx, req.GetTripId(), req.GetMemberId())
	if err != nil {
		return nil, err
	}
	c.Role = req.GetRole()
	return c, s.repo.Save(ctx, c)
}

// Remove - 제거.
func (s *CompanionService) Remove(ctx context.Context, tripID, memberID string) error {
	return s.repo.Delete(ctx, tripID, memberID)
}

// === Tag ===

// TagRepo - 메모리.
type TagRepo struct {
	mu      sync.RWMutex
	tags    map[string]map[string]*journeyv1.Tag // userID -> (tagID -> Tag)
	tripTag map[string]map[string]*journeyv1.Tag // tripID -> (tagID -> denormalized Tag)
}

// NewTagRepo - 생성.
func NewTagRepo() *TagRepo {
	return &TagRepo{
		tags:    make(map[string]map[string]*journeyv1.Tag),
		tripTag: make(map[string]map[string]*journeyv1.Tag),
	}
}

// Save - tag 저장.
func (r *TagRepo) Save(_ context.Context, t *journeyv1.Tag) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tags[t.GetUserId()]; !ok {
		r.tags[t.GetUserId()] = make(map[string]*journeyv1.Tag)
	}
	r.tags[t.GetUserId()][t.GetId()] = t
	return nil
}

// Get - 조회 (userID 필수).
func (r *TagRepo) Get(_ context.Context, userID, id string) (*journeyv1.Tag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if byUser, ok := r.tags[userID]; ok {
		if t, ok := byUser[id]; ok {
			return t, nil
		}
	}
	return nil, ErrNotFound
}

// ListByUser - 사용자 태그 목록.
func (r *TagRepo) ListByUser(_ context.Context, userID string) ([]*journeyv1.Tag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*journeyv1.Tag, 0)
	for _, t := range r.tags[userID] {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].GetName() < out[j].GetName() })
	return out, nil
}

// Delete - 삭제 (userID 필수; 모든 trip 연결도 끊음).
func (r *TagRepo) Delete(_ context.Context, userID, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if byUser, ok := r.tags[userID]; ok {
		delete(byUser, id)
	}
	for tripID, set := range r.tripTag {
		delete(set, id)
		if len(set) == 0 {
			delete(r.tripTag, tripID)
		}
	}
	return nil
}

// Attach - trip 에 tag 부착 (denormalized).
func (r *TagRepo) Attach(_ context.Context, tripID string, tag *journeyv1.Tag) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tripTag[tripID]; !ok {
		r.tripTag[tripID] = make(map[string]*journeyv1.Tag)
	}
	r.tripTag[tripID][tag.GetId()] = tag
	return nil
}

// Detach - 제거.
func (r *TagRepo) Detach(_ context.Context, tripID, tagID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if set, ok := r.tripTag[tripID]; ok {
		delete(set, tagID)
	}
	return nil
}

// ListByTrip - trip 의 tag 목록.
func (r *TagRepo) ListByTrip(_ context.Context, tripID string) ([]*journeyv1.Tag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*journeyv1.Tag, 0)
	if set, ok := r.tripTag[tripID]; ok {
		for _, t := range set {
			out = append(out, t)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].GetName() < out[j].GetName() })
	return out, nil
}

// DetachAllFromTrip - trip 삭제 시.
func (r *TagRepo) DetachAllFromTrip(_ context.Context, tripID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tripTag, tripID)
	return nil
}

// TagService - 로직.
type TagService struct {
	repo TagRepository
}

// NewTagService - 생성.
func NewTagService(repo TagRepository) *TagService { return &TagService{repo: repo} }

// Create - 생성.
func (s *TagService) Create(ctx context.Context, userID string, req *journeyv1.CreateTagRequest) (*journeyv1.Tag, error) {
	t := &journeyv1.Tag{Id: uuid.NewString(), UserId: userID, Name: req.GetName(), Color: req.GetColor()}
	return t, s.repo.Save(ctx, t)
}

// Get - 조회.
func (s *TagService) Get(ctx context.Context, userID, id string) (*journeyv1.Tag, error) {
	return s.repo.Get(ctx, userID, id)
}

// List - 사용자 태그.
func (s *TagService) List(ctx context.Context, userID string) ([]*journeyv1.Tag, error) {
	return s.repo.ListByUser(ctx, userID)
}

// Delete - 삭제.
func (s *TagService) Delete(ctx context.Context, userID, id string) error {
	return s.repo.Delete(ctx, userID, id)
}

// Attach - 부착. 태그 본체를 미리 조회해 denormalized row 를 쓴다.
func (s *TagService) Attach(ctx context.Context, userID, tripID, tagID string) error {
	t, err := s.repo.Get(ctx, userID, tagID)
	if err != nil {
		return err
	}
	return s.repo.Attach(ctx, tripID, t)
}

// Detach - 제거.
func (s *TagService) Detach(ctx context.Context, tripID, tagID string) error {
	return s.repo.Detach(ctx, tripID, tagID)
}

// ListByTrip - trip 의 태그.
func (s *TagService) ListByTrip(ctx context.Context, tripID string) ([]*journeyv1.Tag, error) {
	return s.repo.ListByTrip(ctx, tripID)
}

// === Template ===

// TemplateRepo - 메모리.
type TemplateRepo struct {
	mu    sync.RWMutex
	store map[string]map[string]*journeyv1.Template // userID -> (templateID -> Template)
}

// NewTemplateRepo - 생성.
func NewTemplateRepo() *TemplateRepo {
	return &TemplateRepo{store: make(map[string]map[string]*journeyv1.Template)}
}

// Save - 저장.
func (r *TemplateRepo) Save(_ context.Context, t *journeyv1.Template) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[t.GetUserId()]; !ok {
		r.store[t.GetUserId()] = make(map[string]*journeyv1.Template)
	}
	r.store[t.GetUserId()][t.GetId()] = t
	return nil
}

// Get - 조회.
func (r *TemplateRepo) Get(_ context.Context, userID, id string) (*journeyv1.Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if byUser, ok := r.store[userID]; ok {
		if t, ok := byUser[id]; ok {
			return t, nil
		}
	}
	return nil, ErrNotFound
}

// ListByUser - 목록.
func (r *TemplateRepo) ListByUser(_ context.Context, userID string) ([]*journeyv1.Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*journeyv1.Template, 0)
	for _, t := range r.store[userID] {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].GetName() < out[j].GetName() })
	return out, nil
}

// Delete - 삭제.
func (r *TemplateRepo) Delete(_ context.Context, userID, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if byUser, ok := r.store[userID]; ok {
		delete(byUser, id)
	}
	return nil
}

// TemplateService - 로직.
type TemplateService struct {
	repo TemplateRepository
	now  func() time.Time
}

// NewTemplateService - 생성.
func NewTemplateService(repo TemplateRepository) *TemplateService {
	return &TemplateService{repo: repo, now: time.Now}
}

// Save - 신규 등록.
func (s *TemplateService) Save(ctx context.Context, userID string, req *journeyv1.SaveTripAsTemplateRequest) (*journeyv1.Template, error) {
	now := timestamppb.New(s.now().UTC())
	t := &journeyv1.Template{
		Id:           uuid.NewString(),
		UserId:       userID,
		Name:         req.GetName(),
		SourceTripId: req.GetTripId(),
		Audit:        &journeyv1.AuditTimestamps{CreatedAt: now, UpdatedAt: now},
	}
	return t, s.repo.Save(ctx, t)
}

// Get - 조회.
func (s *TemplateService) Get(ctx context.Context, userID, id string) (*journeyv1.Template, error) {
	return s.repo.Get(ctx, userID, id)
}

// List - 사용자 템플릿.
func (s *TemplateService) List(ctx context.Context, userID string) ([]*journeyv1.Template, error) {
	return s.repo.ListByUser(ctx, userID)
}

// Delete - 삭제.
func (s *TemplateService) Delete(ctx context.Context, userID, id string) error {
	return s.repo.Delete(ctx, userID, id)
}

// === FavoritePlace ===

// FavoriteRepo - 메모리.
type FavoriteRepo struct {
	mu    sync.RWMutex
	store map[string]map[string]*journeyv1.FavoritePlace // userID -> (placeID -> FavoritePlace)
}

// NewFavoriteRepo - 생성.
func NewFavoriteRepo() *FavoriteRepo {
	return &FavoriteRepo{store: make(map[string]map[string]*journeyv1.FavoritePlace)}
}

// Save - 저장.
func (r *FavoriteRepo) Save(_ context.Context, p *journeyv1.FavoritePlace) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[p.GetUserId()]; !ok {
		r.store[p.GetUserId()] = make(map[string]*journeyv1.FavoritePlace)
	}
	r.store[p.GetUserId()][p.GetId()] = p
	return nil
}

// ListByUser - 목록.
func (r *FavoriteRepo) ListByUser(_ context.Context, userID string) ([]*journeyv1.FavoritePlace, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*journeyv1.FavoritePlace, 0)
	for _, p := range r.store[userID] {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].GetName() < out[j].GetName() })
	return out, nil
}

// Delete - 삭제.
func (r *FavoriteRepo) Delete(_ context.Context, userID, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if byUser, ok := r.store[userID]; ok {
		delete(byUser, id)
	}
	return nil
}

// FavoriteService - 로직.
type FavoriteService struct {
	repo FavoriteRepository
	now  func() time.Time
}

// NewFavoriteService - 생성.
func NewFavoriteService(repo FavoriteRepository) *FavoriteService {
	return &FavoriteService{repo: repo, now: time.Now}
}

// Add - 추가.
func (s *FavoriteService) Add(ctx context.Context, userID string, req *journeyv1.AddFavoritePlaceRequest) (*journeyv1.FavoritePlace, error) {
	now := timestamppb.New(s.now().UTC())
	p := &journeyv1.FavoritePlace{
		Id:            uuid.NewString(),
		UserId:        userID,
		Name:          req.GetName(),
		Location:      req.GetLocation(),
		GooglePlaceId: req.GetGooglePlaceId(),
		Memo:          req.GetMemo(),
		Audit:         &journeyv1.AuditTimestamps{CreatedAt: now, UpdatedAt: now},
	}
	return p, s.repo.Save(ctx, p)
}

// List - 목록.
func (s *FavoriteService) List(ctx context.Context, userID string) ([]*journeyv1.FavoritePlace, error) {
	return s.repo.ListByUser(ctx, userID)
}

// Remove - 삭제.
func (s *FavoriteService) Remove(ctx context.Context, userID, id string) error {
	return s.repo.Delete(ctx, userID, id)
}

// === Share ===

// ShareRepo - 메모리.
type ShareRepo struct {
	mu    sync.RWMutex
	store map[string]*journeyv1.Share // code -> Share
}

// NewShareRepo - 생성.
func NewShareRepo() *ShareRepo { return &ShareRepo{store: make(map[string]*journeyv1.Share)} }

// Save - 저장.
func (r *ShareRepo) Save(_ context.Context, s *journeyv1.Share) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[s.GetCode()] = s
	return nil
}

// Get - 조회.
func (r *ShareRepo) Get(_ context.Context, code string) (*journeyv1.Share, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.store[code]
	if !ok {
		return nil, ErrNotFound
	}
	return s, nil
}

// Delete - 삭제.
func (r *ShareRepo) Delete(_ context.Context, code string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[code]; !ok {
		return ErrNotFound
	}
	delete(r.store, code)
	return nil
}

// ShareService - 로직.
type ShareService struct {
	repo ShareRepository
	now  func() time.Time
}

// NewShareService - 생성.
func NewShareService(repo ShareRepository) *ShareService {
	return &ShareService{repo: repo, now: time.Now}
}

func newShareCode() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// Create - 공유 링크 생성.
func (s *ShareService) Create(ctx context.Context, req *journeyv1.CreateShareRequest) (*journeyv1.Share, error) {
	now := timestamppb.New(s.now().UTC())
	sh := &journeyv1.Share{
		Code:       newShareCode(),
		TripId:     req.GetTripId(),
		Permission: req.GetPermission(),
		ExpiresAt:  req.GetExpiresAt(),
		Audit:      &journeyv1.AuditTimestamps{CreatedAt: now, UpdatedAt: now},
	}
	return sh, s.repo.Save(ctx, sh)
}

// Get - 코드 조회. 만료 검증 포함.
func (s *ShareService) Get(ctx context.Context, code string) (*journeyv1.Share, error) {
	sh, err := s.repo.Get(ctx, code)
	if err != nil {
		return nil, err
	}
	if sh.GetExpiresAt() != nil && sh.GetExpiresAt().AsTime().Before(s.now().UTC()) {
		return nil, ErrNotFound
	}
	return sh, nil
}

// Revoke - 무효화.
func (s *ShareService) Revoke(ctx context.Context, code string) error {
	return s.repo.Delete(ctx, code)
}
