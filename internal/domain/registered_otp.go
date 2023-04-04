package domain

import (
	"time"

	"gorm.io/gorm"
)

// RegisteredOTP object for table `registered_otp`.
type RegisteredOTP struct {
	ID        uint `gorm:"primaryKey"`
	Code      string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (o *RegisteredOTP) TableName() string {
	return "registered_otp"
}
