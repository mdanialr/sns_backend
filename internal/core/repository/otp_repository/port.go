package otp_repository

import (
	"context"

	"github.com/mdanialr/sns_backend/internal/domain"
)

// IRepository an interface that should be used when dealing with object
// domain.RegisteredOTP.
type IRepository interface {
	// GetByCode retrieve a domain.RegisteredOTP by the given code, also return
	// error if any including record not found.
	GetByCode(ctx context.Context, code string) (*domain.RegisteredOTP, error)
	// Create save new instance of domain.RegisteredOTP that's only need given
	// code.
	Create(ctx context.Context, code string) (*domain.RegisteredOTP, error)
	// DeleteAll batch delete all records of domain.RegisteredOTP.
	DeleteAll(ctx context.Context) error
}
