package sns_repository

import (
	"context"

	r "github.com/mdanialr/sns_backend/internal/core/repository"
	"github.com/mdanialr/sns_backend/internal/domain"
)

// IRepository an interface that may be used when dealing with object
// domain.SNS.
type IRepository interface {
	// FindShorten retrieve all shorten data.
	FindShorten(ctx context.Context, opts ...r.IOptions) ([]*domain.SNS, error)
	// FindSend retrieve all send data.
	FindSend(ctx context.Context, opts ...r.IOptions) ([]*domain.SNS, error)
	// GetByID retrieve a domain.SNS by given id and optionally select which
	// columns to be retrieved. Returned domain.SNS should be nil even if there
	// is any error.
	GetByID(ctx context.Context, id uint, opts ...r.IOptions) (*domain.SNS, error)
	// GetByUrl same as GetByID but use the url field instead. Returned
	// domain.SNS should be nil even if there is any error.
	GetByUrl(ctx context.Context, url string, opts ...r.IOptions) (*domain.SNS, error)
	// Create save given sns. Return the newly saved object that's the primary
	// key should be filled already.
	Create(ctx context.Context, sns *domain.SNS) (*domain.SNS, error)
	// Update do update given sns using given cons if any. Default is using
	// provided primary key in sns param as the conditions.
	Update(ctx context.Context, sns *domain.SNS, opts ...r.IOptions) (*domain.SNS, error)
	// DeleteByID delete an object that's has given id as their primary key.
	DeleteByID(ctx context.Context, id uint) error
}
