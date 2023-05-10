package domain

import (
	"time"

	"gorm.io/gorm"
)

// SNS object for table `sns`.
type SNS struct {
	ID          uint `gorm:"primaryKey"`
	Url         string
	Shorten     *string
	Send        *string
	FileSize    *string
	IsPermanent *bool
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	DeletedAt   gorm.DeletedAt ` gorm:"index"`
}

func (s *SNS) TableName() string {
	return "sns"
}
