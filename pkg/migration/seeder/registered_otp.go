package seeder

import (
	"github.com/mdanialr/sns_backend/internal/domain"
	"gorm.io/gorm"
)

func registeredOTP(db *gorm.DB) {
	var samples = []domain.RegisteredOTP{
		{Code: "123456"},
		{Code: "432182"},
		{Code: "322012"},
	}
	for _, sample := range samples {
		db.Create(&sample)
	}
}
