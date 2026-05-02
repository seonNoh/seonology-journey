// Package record - Expense / Note / Checklist / Reservation 도메인.
package record

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
var ErrNotFound = errors.New("record: not found")

// === Expense ===

// ExpenseRepo - 메모리.
type ExpenseRepo struct {
	mu    sync.RWMutex
	store map[string]*journeyv1.Expense
}

// NewExpenseRepo - 생성.
func NewExpenseRepo() *ExpenseRepo { return &ExpenseRepo{store: make(map[string]*journeyv1.Expense)} }

// Get - id 조회.
func (r *ExpenseRepo) Get(_ context.Context, id string) (*journeyv1.Expense, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.store[id]
	if !ok {
		return nil, ErrNotFound
	}
	return e, nil
}

// Save - 저장.
func (r *ExpenseRepo) Save(_ context.Context, e *journeyv1.Expense) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[e.GetId()] = e
	return nil
}

// Delete - 삭제.
func (r *ExpenseRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[id]; !ok {
		return ErrNotFound
	}
	delete(r.store, id)
	return nil
}

// ListByTrip - trip 의 지출 목록 (spent_at 오름차순).
func (r *ExpenseRepo) ListByTrip(_ context.Context, tripID string) ([]*journeyv1.Expense, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*journeyv1.Expense, 0)
	for _, e := range r.store {
		if e.GetTripId() == tripID {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].GetSpentAt().AsTime().Before(out[j].GetSpentAt().AsTime())
	})
	return out, nil
}

// DeleteByTrip - cascade.
func (r *ExpenseRepo) DeleteByTrip(_ context.Context, tripID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for k, e := range r.store {
		if e.GetTripId() == tripID {
			delete(r.store, k)
		}
	}
	return nil
}

// ExpenseService - 로직.
type ExpenseService struct {
	repo *ExpenseRepo
	now  func() time.Time
}

// NewExpenseService - 생성.
func NewExpenseService(repo *ExpenseRepo) *ExpenseService {
	return &ExpenseService{repo: repo, now: time.Now}
}

// Create - 생성.
func (s *ExpenseService) Create(ctx context.Context, req *journeyv1.CreateExpenseRequest) (*journeyv1.Expense, error) {
	now := timestamppb.New(s.now().UTC())
	e := &journeyv1.Expense{
		Id:            uuid.NewString(),
		TripId:        req.GetTripId(),
		DayId:         req.GetDayId(),
		Category:      req.GetCategory(),
		Amount:        req.GetAmount(),
		PaymentMethod: req.GetPaymentMethod(),
		Description:   req.GetDescription(),
		SpentAt:       req.GetSpentAt(),
		Audit:         &journeyv1.AuditTimestamps{CreatedAt: now, UpdatedAt: now},
	}
	if e.SpentAt == nil {
		e.SpentAt = now
	}
	return e, s.repo.Save(ctx, e)
}

// Get - 조회.
func (s *ExpenseService) Get(ctx context.Context, id string) (*journeyv1.Expense, error) {
	return s.repo.Get(ctx, id)
}

// List - trip 의 지출.
func (s *ExpenseService) List(ctx context.Context, tripID string) ([]*journeyv1.Expense, error) {
	return s.repo.ListByTrip(ctx, tripID)
}

// Update - 갱신.
func (s *ExpenseService) Update(ctx context.Context, req *journeyv1.UpdateExpenseRequest) (*journeyv1.Expense, error) {
	e, err := s.repo.Get(ctx, req.GetExpenseId())
	if err != nil {
		return nil, err
	}
	if req.GetCategory() != journeyv1.ExpenseCategory_EXPENSE_CATEGORY_UNSPECIFIED {
		e.Category = req.GetCategory()
	}
	if req.GetAmount() != nil {
		e.Amount = req.GetAmount()
	}
	if req.GetPaymentMethod() != journeyv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		e.PaymentMethod = req.GetPaymentMethod()
	}
	if req.GetDescription() != "" {
		e.Description = req.GetDescription()
	}
	if req.GetDayId() != "" {
		e.DayId = req.GetDayId()
	}
	if e.Audit == nil {
		e.Audit = &journeyv1.AuditTimestamps{}
	}
	e.Audit.UpdatedAt = timestamppb.New(s.now().UTC())
	return e, s.repo.Save(ctx, e)
}

// Delete - 삭제.
func (s *ExpenseService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// Summary - trip 통계.
func (s *ExpenseService) Summary(ctx context.Context, tripID string) (*journeyv1.ExpenseSummary, error) {
	list, err := s.repo.ListByTrip(ctx, tripID)
	if err != nil {
		return nil, err
	}
	sum := &journeyv1.ExpenseSummary{TripId: tripID}
	if len(list) == 0 {
		return sum, nil
	}
	cur := list[0].GetAmount().GetCurrency()
	grand := int64(0)
	byCat := map[journeyv1.ExpenseCategory]int64{}
	byDay := map[string]int64{}
	for _, e := range list {
		amt := e.GetAmount().GetAmount()
		grand += amt
		byCat[e.GetCategory()] += amt
		date := e.GetSpentAt().AsTime().UTC().Format("2006-01-02")
		byDay[date] += amt
	}
	sum.GrandTotal = &journeyv1.Money{Currency: cur, Amount: grand}
	for c, t := range byCat {
		sum.ByCategory = append(sum.ByCategory, &journeyv1.ExpenseSummary_ByCategory{
			Category: c, Total: &journeyv1.Money{Currency: cur, Amount: t},
		})
	}
	dates := make([]string, 0, len(byDay))
	for d := range byDay {
		dates = append(dates, d)
	}
	sort.Strings(dates)
	for _, d := range dates {
		sum.ByDay = append(sum.ByDay, &journeyv1.ExpenseSummary_ByDay{
			Date: d, Total: &journeyv1.Money{Currency: cur, Amount: byDay[d]},
		})
	}
	return sum, nil
}

// === Note ===

// NoteRepo - 메모리.
type NoteRepo struct {
	mu    sync.RWMutex
	store map[string]*journeyv1.Note
}

// NewNoteRepo - 생성.
func NewNoteRepo() *NoteRepo { return &NoteRepo{store: make(map[string]*journeyv1.Note)} }

// Get - 조회.
func (r *NoteRepo) Get(_ context.Context, id string) (*journeyv1.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	n, ok := r.store[id]
	if !ok {
		return nil, ErrNotFound
	}
	return n, nil
}

// Save - 저장.
func (r *NoteRepo) Save(_ context.Context, n *journeyv1.Note) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[n.GetId()] = n
	return nil
}

// Delete - 삭제.
func (r *NoteRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[id]; !ok {
		return ErrNotFound
	}
	delete(r.store, id)
	return nil
}

// ListByTrip - 목록.
func (r *NoteRepo) ListByTrip(_ context.Context, tripID string) ([]*journeyv1.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*journeyv1.Note, 0)
	for _, n := range r.store {
		if n.GetTripId() == tripID {
			out = append(out, n)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].GetAudit().GetCreatedAt().AsTime().Before(out[j].GetAudit().GetCreatedAt().AsTime())
	})
	return out, nil
}

// DeleteByTrip - cascade.
func (r *NoteRepo) DeleteByTrip(_ context.Context, tripID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for k, n := range r.store {
		if n.GetTripId() == tripID {
			delete(r.store, k)
		}
	}
	return nil
}

// NoteService - 로직.
type NoteService struct {
	repo *NoteRepo
	now  func() time.Time
}

// NewNoteService - 생성.
func NewNoteService(repo *NoteRepo) *NoteService {
	return &NoteService{repo: repo, now: time.Now}
}

// Create - 생성.
func (s *NoteService) Create(ctx context.Context, req *journeyv1.CreateNoteRequest) (*journeyv1.Note, error) {
	now := timestamppb.New(s.now().UTC())
	n := &journeyv1.Note{
		Id:      uuid.NewString(),
		TripId:  req.GetTripId(),
		DayId:   req.GetDayId(),
		Content: req.GetContent(),
		Mood:    req.GetMood(),
		Audit:   &journeyv1.AuditTimestamps{CreatedAt: now, UpdatedAt: now},
	}
	return n, s.repo.Save(ctx, n)
}

// Get - 조회.
func (s *NoteService) Get(ctx context.Context, id string) (*journeyv1.Note, error) {
	return s.repo.Get(ctx, id)
}

// List - trip 의 메모.
func (s *NoteService) List(ctx context.Context, tripID string) ([]*journeyv1.Note, error) {
	return s.repo.ListByTrip(ctx, tripID)
}

// Update - 갱신.
func (s *NoteService) Update(ctx context.Context, req *journeyv1.UpdateNoteRequest) (*journeyv1.Note, error) {
	n, err := s.repo.Get(ctx, req.GetNoteId())
	if err != nil {
		return nil, err
	}
	if req.GetContent() != "" {
		n.Content = req.GetContent()
	}
	if req.GetMood() != "" {
		n.Mood = req.GetMood()
	}
	if n.Audit == nil {
		n.Audit = &journeyv1.AuditTimestamps{}
	}
	n.Audit.UpdatedAt = timestamppb.New(s.now().UTC())
	return n, s.repo.Save(ctx, n)
}

// Delete - 삭제.
func (s *NoteService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// === Checklist ===

// ChecklistRepo - 메모리.
type ChecklistRepo struct {
	mu    sync.RWMutex
	store map[string]*journeyv1.ChecklistItem
}

// NewChecklistRepo - 생성.
func NewChecklistRepo() *ChecklistRepo {
	return &ChecklistRepo{store: make(map[string]*journeyv1.ChecklistItem)}
}

// Get - 조회.
func (r *ChecklistRepo) Get(_ context.Context, id string) (*journeyv1.ChecklistItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.store[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

// Save - 저장.
func (r *ChecklistRepo) Save(_ context.Context, c *journeyv1.ChecklistItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[c.GetId()] = c
	return nil
}

// Delete - 삭제.
func (r *ChecklistRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[id]; !ok {
		return ErrNotFound
	}
	delete(r.store, id)
	return nil
}

// ListByTrip - 목록.
func (r *ChecklistRepo) ListByTrip(_ context.Context, tripID string) ([]*journeyv1.ChecklistItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*journeyv1.ChecklistItem, 0)
	for _, c := range r.store {
		if c.GetTripId() == tripID {
			out = append(out, c)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].GetCategory() != out[j].GetCategory() {
			return out[i].GetCategory() < out[j].GetCategory()
		}
		return out[i].GetAudit().GetCreatedAt().AsTime().Before(out[j].GetAudit().GetCreatedAt().AsTime())
	})
	return out, nil
}

// DeleteByTrip - cascade.
func (r *ChecklistRepo) DeleteByTrip(_ context.Context, tripID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for k, c := range r.store {
		if c.GetTripId() == tripID {
			delete(r.store, k)
		}
	}
	return nil
}

// ChecklistService - 로직.
type ChecklistService struct {
	repo *ChecklistRepo
	now  func() time.Time
}

// NewChecklistService - 생성.
func NewChecklistService(repo *ChecklistRepo) *ChecklistService {
	return &ChecklistService{repo: repo, now: time.Now}
}

// Create - 생성.
func (s *ChecklistService) Create(ctx context.Context, req *journeyv1.CreateChecklistItemRequest) (*journeyv1.ChecklistItem, error) {
	now := timestamppb.New(s.now().UTC())
	c := &journeyv1.ChecklistItem{
		Id:       uuid.NewString(),
		TripId:   req.GetTripId(),
		Category: req.GetCategory(),
		Item:     req.GetItem(),
		Audit:    &journeyv1.AuditTimestamps{CreatedAt: now, UpdatedAt: now},
	}
	return c, s.repo.Save(ctx, c)
}

// Get - 조회.
func (s *ChecklistService) Get(ctx context.Context, id string) (*journeyv1.ChecklistItem, error) {
	return s.repo.Get(ctx, id)
}

// List - trip 의 체크리스트.
func (s *ChecklistService) List(ctx context.Context, tripID string) ([]*journeyv1.ChecklistItem, error) {
	return s.repo.ListByTrip(ctx, tripID)
}

// Update - 갱신.
func (s *ChecklistService) Update(ctx context.Context, req *journeyv1.UpdateChecklistItemRequest) (*journeyv1.ChecklistItem, error) {
	c, err := s.repo.Get(ctx, req.GetItemId())
	if err != nil {
		return nil, err
	}
	if req.GetCategory() != journeyv1.ChecklistCategory_CHECKLIST_CATEGORY_UNSPECIFIED {
		c.Category = req.GetCategory()
	}
	if req.GetItem() != "" {
		c.Item = req.GetItem()
	}
	c.IsChecked = req.GetIsChecked()
	if c.Audit == nil {
		c.Audit = &journeyv1.AuditTimestamps{}
	}
	c.Audit.UpdatedAt = timestamppb.New(s.now().UTC())
	return c, s.repo.Save(ctx, c)
}

// Delete - 삭제.
func (s *ChecklistService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// === Reservation ===

// ReservationRepo - 메모리.
type ReservationRepo struct {
	mu    sync.RWMutex
	store map[string]*journeyv1.Reservation
}

// NewReservationRepo - 생성.
func NewReservationRepo() *ReservationRepo {
	return &ReservationRepo{store: make(map[string]*journeyv1.Reservation)}
}

// Get - 조회.
func (r *ReservationRepo) Get(_ context.Context, id string) (*journeyv1.Reservation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.store[id]
	if !ok {
		return nil, ErrNotFound
	}
	return v, nil
}

// Save - 저장.
func (r *ReservationRepo) Save(_ context.Context, v *journeyv1.Reservation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[v.GetId()] = v
	return nil
}

// Delete - 삭제.
func (r *ReservationRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[id]; !ok {
		return ErrNotFound
	}
	delete(r.store, id)
	return nil
}

// ListByTrip - 목록.
func (r *ReservationRepo) ListByTrip(_ context.Context, tripID string) ([]*journeyv1.Reservation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*journeyv1.Reservation, 0)
	for _, v := range r.store {
		if v.GetTripId() == tripID {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].GetReservedAt().AsTime().Before(out[j].GetReservedAt().AsTime())
	})
	return out, nil
}

// DeleteByTrip - cascade.
func (r *ReservationRepo) DeleteByTrip(_ context.Context, tripID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for k, v := range r.store {
		if v.GetTripId() == tripID {
			delete(r.store, k)
		}
	}
	return nil
}

// ReservationService - 로직.
type ReservationService struct {
	repo *ReservationRepo
	now  func() time.Time
}

// NewReservationService - 생성.
func NewReservationService(repo *ReservationRepo) *ReservationService {
	return &ReservationService{repo: repo, now: time.Now}
}

// Create - 생성.
func (s *ReservationService) Create(ctx context.Context, req *journeyv1.CreateReservationRequest) (*journeyv1.Reservation, error) {
	now := timestamppb.New(s.now().UTC())
	v := &journeyv1.Reservation{
		Id:              uuid.NewString(),
		TripId:          req.GetTripId(),
		Type:            req.GetType(),
		Vendor:          req.GetVendor(),
		ConfirmNumber:   req.GetConfirmNumber(),
		ReservedAt:      req.GetReservedAt(),
		Cost:            req.GetCost(),
		AttachmentS3Key: req.GetAttachmentS3Key(),
		Notes:           req.GetNotes(),
		Audit:           &journeyv1.AuditTimestamps{CreatedAt: now, UpdatedAt: now},
	}
	return v, s.repo.Save(ctx, v)
}

// Get - 조회.
func (s *ReservationService) Get(ctx context.Context, id string) (*journeyv1.Reservation, error) {
	return s.repo.Get(ctx, id)
}

// List - trip 의 예약.
func (s *ReservationService) List(ctx context.Context, tripID string) ([]*journeyv1.Reservation, error) {
	return s.repo.ListByTrip(ctx, tripID)
}

// Update - 갱신.
func (s *ReservationService) Update(ctx context.Context, req *journeyv1.UpdateReservationRequest) (*journeyv1.Reservation, error) {
	v, err := s.repo.Get(ctx, req.GetReservationId())
	if err != nil {
		return nil, err
	}
	if req.GetType() != journeyv1.ReservationType_RESERVATION_TYPE_UNSPECIFIED {
		v.Type = req.GetType()
	}
	if req.GetVendor() != "" {
		v.Vendor = req.GetVendor()
	}
	if req.GetConfirmNumber() != "" {
		v.ConfirmNumber = req.GetConfirmNumber()
	}
	if req.GetReservedAt() != nil {
		v.ReservedAt = req.GetReservedAt()
	}
	if req.GetCost() != nil {
		v.Cost = req.GetCost()
	}
	if req.GetAttachmentS3Key() != "" {
		v.AttachmentS3Key = req.GetAttachmentS3Key()
	}
	if req.GetNotes() != "" {
		v.Notes = req.GetNotes()
	}
	if v.Audit == nil {
		v.Audit = &journeyv1.AuditTimestamps{}
	}
	v.Audit.UpdatedAt = timestamppb.New(s.now().UTC())
	return v, s.repo.Save(ctx, v)
}

// Delete - 삭제.
func (s *ReservationService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
