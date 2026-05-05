package social

import (
	"context"

	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
)

// CompanionRepository defines the interface for companion persistence.
type CompanionRepository interface {
	Save(ctx context.Context, c *journeyv1.Companion) error
	Get(ctx context.Context, tripID, memberID string) (*journeyv1.Companion, error)
	ListByTrip(ctx context.Context, tripID string) ([]*journeyv1.Companion, error)
	Delete(ctx context.Context, tripID, memberID string) error
	DeleteByTrip(ctx context.Context, tripID string) error
}

// TagRepository defines the interface for tag persistence.
type TagRepository interface {
	Save(ctx context.Context, t *journeyv1.Tag) error
	Get(ctx context.Context, id string) (*journeyv1.Tag, error)
	ListByUser(ctx context.Context, userID string) ([]*journeyv1.Tag, error)
	Delete(ctx context.Context, id string) error
	Attach(ctx context.Context, tripID, tagID string) error
	Detach(ctx context.Context, tripID, tagID string) error
	ListByTrip(ctx context.Context, tripID string) ([]*journeyv1.Tag, error)
	DetachAllFromTrip(ctx context.Context, tripID string) error
}

// TemplateRepository defines the interface for template persistence.
type TemplateRepository interface {
	Save(ctx context.Context, t *journeyv1.Template) error
	Get(ctx context.Context, id string) (*journeyv1.Template, error)
	ListByUser(ctx context.Context, userID string) ([]*journeyv1.Template, error)
	Delete(ctx context.Context, id string) error
}

// FavoriteRepository defines the interface for favorite place persistence.
type FavoriteRepository interface {
	Save(ctx context.Context, p *journeyv1.FavoritePlace) error
	ListByUser(ctx context.Context, userID string) ([]*journeyv1.FavoritePlace, error)
	Delete(ctx context.Context, id string) error
}

// ShareRepository defines the interface for share persistence.
type ShareRepository interface {
	Save(ctx context.Context, s *journeyv1.Share) error
	Get(ctx context.Context, code string) (*journeyv1.Share, error)
	Delete(ctx context.Context, code string) error
}
