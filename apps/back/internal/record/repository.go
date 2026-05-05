package record

import (
	"context"

	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
)

// ExpenseRepository defines the interface for expense persistence.
type ExpenseRepository interface {
	Get(ctx context.Context, id string) (*journeyv1.Expense, error)
	Save(ctx context.Context, e *journeyv1.Expense) error
	Delete(ctx context.Context, id string) error
	ListByTrip(ctx context.Context, tripID string) ([]*journeyv1.Expense, error)
	DeleteByTrip(ctx context.Context, tripID string) error
}

// NoteRepository defines the interface for note persistence.
type NoteRepository interface {
	Get(ctx context.Context, id string) (*journeyv1.Note, error)
	Save(ctx context.Context, n *journeyv1.Note) error
	Delete(ctx context.Context, id string) error
	ListByTrip(ctx context.Context, tripID string) ([]*journeyv1.Note, error)
	DeleteByTrip(ctx context.Context, tripID string) error
}

// ChecklistRepository defines the interface for checklist persistence.
type ChecklistRepository interface {
	Get(ctx context.Context, id string) (*journeyv1.ChecklistItem, error)
	Save(ctx context.Context, c *journeyv1.ChecklistItem) error
	Delete(ctx context.Context, id string) error
	ListByTrip(ctx context.Context, tripID string) ([]*journeyv1.ChecklistItem, error)
	DeleteByTrip(ctx context.Context, tripID string) error
}

// ReservationRepository defines the interface for reservation persistence.
type ReservationRepository interface {
	Get(ctx context.Context, id string) (*journeyv1.Reservation, error)
	Save(ctx context.Context, v *journeyv1.Reservation) error
	Delete(ctx context.Context, id string) error
	ListByTrip(ctx context.Context, tripID string) ([]*journeyv1.Reservation, error)
	DeleteByTrip(ctx context.Context, tripID string) error
}
