package shorten_service

import (
	"context"

	req "github.com/mdanialr/sns_backend/internal/requests"
	res "github.com/mdanialr/sns_backend/internal/responses"
)

// IService an interface that should be used when dealing with sns.
type IService interface {
	// Index retrieve all shorten data with a pagination if provided from
	// request query params.
	Index(context.Context, *req.Shorten) (*res.ShortenIndexResponse, error)
	// Create save a new Shorten instance based on given request. Return the
	// newly created Shorten back along with error if any.
	Create(context.Context, *req.Shorten) (*res.ShortenResponse, error)
	// Update do update an existing Shorten instance based on ID in given
	// request.
	Update(context.Context, *req.ShortenUpdate) (*res.ShortenResponse, error)
	// Delete remove an SNS data from DB using given id as the condition.
	Delete(context.Context, *req.ShortenDelete) error
}
