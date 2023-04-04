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
}
