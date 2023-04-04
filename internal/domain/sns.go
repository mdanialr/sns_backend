package domain

import (
	"time"

	"gorm.io/gorm"
)

// SNS object for table `sns`.
type SNS struct {
	ID          uint           `json:"id,omitempty" gorm:"primaryKey"`
	Url         string         `json:"url,omitempty"`
	Shorten     *string        `json:"shorten,omitempty"`
	Send        *string        `json:"send,omitempty"`
	FileSize    *string        `json:"file_size,omitempty"`
	IsPermanent *bool          `json:"is_permanent,omitempty"`
	CreatedAt   *time.Time     `json:"created_at,omitempty"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (s *SNS) TableName() string {
	return "sns"
}
