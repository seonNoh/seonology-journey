package server

import (
	"context"

	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
)

// requireTripOwner - trip 소유 검증.
func (s *JourneyServer) requireTripOwner(ctx context.Context, owner, tripID string) error {
	if _, err := s.d.Trip.Get(ctx, owner, tripID); err != nil {
		return mapErr(err)
	}
	return nil
}

// === Expense ===

// CreateExpense implements JourneyService.
func (s *JourneyServer) CreateExpense(ctx context.Context, req *journeyv1.CreateExpenseRequest) (*journeyv1.CreateExpenseResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.requireTripOwner(ctx, owner, req.GetTripId()); err != nil {
		return nil, err
	}
	e, err := s.d.Expense.Create(ctx, req)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.CreateExpenseResponse{Expense: e}, nil
}

// ListExpenses implements JourneyService.
func (s *JourneyServer) ListExpenses(ctx context.Context, req *journeyv1.ListExpensesRequest) (*journeyv1.ListExpensesResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.requireTripOwner(ctx, owner, req.GetTripId()); err != nil {
		return nil, err
	}
	out, err := s.d.Expense.List(ctx, req.GetTripId())
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.ListExpensesResponse{Expenses: out, Page: &journeyv1.PageInfo{HasMore: false}}, nil
}

// UpdateExpense implements JourneyService.
func (s *JourneyServer) UpdateExpense(ctx context.Context, req *journeyv1.UpdateExpenseRequest) (*journeyv1.UpdateExpenseResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	prev, err := s.d.Expense.Get(ctx, req.GetExpenseId())
	if err != nil {
		return nil, mapErr(err)
	}
	if err := s.requireTripOwner(ctx, owner, prev.GetTripId()); err != nil {
		return nil, err
	}
	e, err := s.d.Expense.Update(ctx, req)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.UpdateExpenseResponse{Expense: e}, nil
}

// DeleteExpense implements JourneyService.
func (s *JourneyServer) DeleteExpense(ctx context.Context, req *journeyv1.DeleteExpenseRequest) (*journeyv1.DeleteExpenseResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	prev, err := s.d.Expense.Get(ctx, req.GetExpenseId())
	if err != nil {
		return nil, mapErr(err)
	}
	if err := s.requireTripOwner(ctx, owner, prev.GetTripId()); err != nil {
		return nil, err
	}
	if err := s.d.Expense.Delete(ctx, req.GetExpenseId()); err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.DeleteExpenseResponse{}, nil
}

// GetExpenseSummary implements JourneyService.
func (s *JourneyServer) GetExpenseSummary(ctx context.Context, req *journeyv1.GetExpenseSummaryRequest) (*journeyv1.GetExpenseSummaryResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.requireTripOwner(ctx, owner, req.GetTripId()); err != nil {
		return nil, err
	}
	sum, err := s.d.Expense.Summary(ctx, req.GetTripId())
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.GetExpenseSummaryResponse{Summary: sum}, nil
}

// === Note ===

// CreateNote implements JourneyService.
func (s *JourneyServer) CreateNote(ctx context.Context, req *journeyv1.CreateNoteRequest) (*journeyv1.CreateNoteResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.requireTripOwner(ctx, owner, req.GetTripId()); err != nil {
		return nil, err
	}
	n, err := s.d.Note.Create(ctx, req)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.CreateNoteResponse{Note: n}, nil
}

// ListNotes implements JourneyService.
func (s *JourneyServer) ListNotes(ctx context.Context, req *journeyv1.ListNotesRequest) (*journeyv1.ListNotesResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.requireTripOwner(ctx, owner, req.GetTripId()); err != nil {
		return nil, err
	}
	out, err := s.d.Note.List(ctx, req.GetTripId())
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.ListNotesResponse{Notes: out}, nil
}

// UpdateNote implements JourneyService.
func (s *JourneyServer) UpdateNote(ctx context.Context, req *journeyv1.UpdateNoteRequest) (*journeyv1.UpdateNoteResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	prev, err := s.d.Note.Get(ctx, req.GetNoteId())
	if err != nil {
		return nil, mapErr(err)
	}
	if err := s.requireTripOwner(ctx, owner, prev.GetTripId()); err != nil {
		return nil, err
	}
	n, err := s.d.Note.Update(ctx, req)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.UpdateNoteResponse{Note: n}, nil
}

// DeleteNote implements JourneyService.
func (s *JourneyServer) DeleteNote(ctx context.Context, req *journeyv1.DeleteNoteRequest) (*journeyv1.DeleteNoteResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	prev, err := s.d.Note.Get(ctx, req.GetNoteId())
	if err != nil {
		return nil, mapErr(err)
	}
	if err := s.requireTripOwner(ctx, owner, prev.GetTripId()); err != nil {
		return nil, err
	}
	if err := s.d.Note.Delete(ctx, req.GetNoteId()); err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.DeleteNoteResponse{}, nil
}

// === Checklist ===

// CreateChecklistItem implements JourneyService.
func (s *JourneyServer) CreateChecklistItem(ctx context.Context, req *journeyv1.CreateChecklistItemRequest) (*journeyv1.CreateChecklistItemResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.requireTripOwner(ctx, owner, req.GetTripId()); err != nil {
		return nil, err
	}
	c, err := s.d.Checklist.Create(ctx, req)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.CreateChecklistItemResponse{Item: c}, nil
}

// ListChecklistItems implements JourneyService.
func (s *JourneyServer) ListChecklistItems(ctx context.Context, req *journeyv1.ListChecklistItemsRequest) (*journeyv1.ListChecklistItemsResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.requireTripOwner(ctx, owner, req.GetTripId()); err != nil {
		return nil, err
	}
	out, err := s.d.Checklist.List(ctx, req.GetTripId())
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.ListChecklistItemsResponse{Items: out}, nil
}

// UpdateChecklistItem implements JourneyService.
func (s *JourneyServer) UpdateChecklistItem(ctx context.Context, req *journeyv1.UpdateChecklistItemRequest) (*journeyv1.UpdateChecklistItemResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	prev, err := s.d.Checklist.Get(ctx, req.GetItemId())
	if err != nil {
		return nil, mapErr(err)
	}
	if err := s.requireTripOwner(ctx, owner, prev.GetTripId()); err != nil {
		return nil, err
	}
	c, err := s.d.Checklist.Update(ctx, req)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.UpdateChecklistItemResponse{Item: c}, nil
}

// DeleteChecklistItem implements JourneyService.
func (s *JourneyServer) DeleteChecklistItem(ctx context.Context, req *journeyv1.DeleteChecklistItemRequest) (*journeyv1.DeleteChecklistItemResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	prev, err := s.d.Checklist.Get(ctx, req.GetItemId())
	if err != nil {
		return nil, mapErr(err)
	}
	if err := s.requireTripOwner(ctx, owner, prev.GetTripId()); err != nil {
		return nil, err
	}
	if err := s.d.Checklist.Delete(ctx, req.GetItemId()); err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.DeleteChecklistItemResponse{}, nil
}

// === Reservation ===

// CreateReservation implements JourneyService.
func (s *JourneyServer) CreateReservation(ctx context.Context, req *journeyv1.CreateReservationRequest) (*journeyv1.CreateReservationResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.requireTripOwner(ctx, owner, req.GetTripId()); err != nil {
		return nil, err
	}
	r, err := s.d.Reservation.Create(ctx, req)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.CreateReservationResponse{Reservation: r}, nil
}

// ListReservations implements JourneyService.
func (s *JourneyServer) ListReservations(ctx context.Context, req *journeyv1.ListReservationsRequest) (*journeyv1.ListReservationsResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.requireTripOwner(ctx, owner, req.GetTripId()); err != nil {
		return nil, err
	}
	out, err := s.d.Reservation.List(ctx, req.GetTripId())
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.ListReservationsResponse{Reservations: out}, nil
}

// UpdateReservation implements JourneyService.
func (s *JourneyServer) UpdateReservation(ctx context.Context, req *journeyv1.UpdateReservationRequest) (*journeyv1.UpdateReservationResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	prev, err := s.d.Reservation.Get(ctx, req.GetReservationId())
	if err != nil {
		return nil, mapErr(err)
	}
	if err := s.requireTripOwner(ctx, owner, prev.GetTripId()); err != nil {
		return nil, err
	}
	r, err := s.d.Reservation.Update(ctx, req)
	if err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.UpdateReservationResponse{Reservation: r}, nil
}

// DeleteReservation implements JourneyService.
func (s *JourneyServer) DeleteReservation(ctx context.Context, req *journeyv1.DeleteReservationRequest) (*journeyv1.DeleteReservationResponse, error) {
	owner, err := ownerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	prev, err := s.d.Reservation.Get(ctx, req.GetReservationId())
	if err != nil {
		return nil, mapErr(err)
	}
	if err := s.requireTripOwner(ctx, owner, prev.GetTripId()); err != nil {
		return nil, err
	}
	if err := s.d.Reservation.Delete(ctx, req.GetReservationId()); err != nil {
		return nil, mapErr(err)
	}
	return &journeyv1.DeleteReservationResponse{}, nil
}
