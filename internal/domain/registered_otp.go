package domain

import (
	"time"

	"gorm.io/gorm"
)

// RegisteredOTP object for table `registered_otp`.
type RegisteredOTP struct {
	ID        uint           `json:"id,omitempty" gorm:"primaryKey"`
	Code      string         `json:"code,omitempty"`
	CreatedAt *time.Time     `json:"created_at,omitempty"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (o *RegisteredOTP) TableName() string {
	return "registered_otp"
}
