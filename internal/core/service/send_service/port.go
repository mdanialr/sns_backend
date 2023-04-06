package send_service

import (
	"context"

	req "github.com/mdanialr/sns_backend/internal/requests"
	res "github.com/mdanialr/sns_backend/internal/responses"
)

type IService interface {
	// Index retrieve all send data with a pagination if provided from
	// request query params.
	Index(context.Context, *req.Send) (*res.SendIndexResponse, error)
	// Create save a new Send instance based on given request. Return the newly
	// created Send back along with error if any.
	Create(ctx context.Context, req *req.Send) (*res.SendResponse, error)
	// Update do update an existing Send instance based on ID in given
	// request. Return the recently updated Send back along with error if
	// any.
	Update(context.Context, *req.SendUpdate) (*res.SendResponse, error)
	// Delete remove an SNS data from DB using given id as the condition.
	Delete(ctx context.Context, req *req.SendDelete) error
}
