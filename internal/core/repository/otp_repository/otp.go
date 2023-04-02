package otp_repository

import (
	"context"

	"github.com/mdanialr/sns_backend/internal/domain"
	"gorm.io/gorm"
)

type otpRepo struct {
	db *gorm.DB
}

// NewOTPRepository return implementation that can be used to interact with
// object domain.RegisteredOTP.
func NewOTPRepository(db *gorm.DB) IRepository {
	return &otpRepo{db}
}

func (o *otpRepo) GetByCode(ctx context.Context, code string) (*domain.RegisteredOTP, error) {
	ro := domain.RegisteredOTP{Code: code}
	return &ro, o.db.WithContext(ctx).Where(&ro).Select("id").First(&ro).Error
}

func (o *otpRepo) Create(ctx context.Context, code string) (*domain.RegisteredOTP, error) {
	ro := domain.RegisteredOTP{Code: code}
	return &ro, o.db.WithContext(ctx).Create(&ro).Error
}

func (o *otpRepo) DeleteAll(ctx context.Context) error {
	return o.db.WithContext(ctx).
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&domain.RegisteredOTP{}).Error
}
