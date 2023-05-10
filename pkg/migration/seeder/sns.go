package seeder

import (
	"github.com/mdanialr/sns_backend/internal/domain"
	h "github.com/mdanialr/sns_backend/pkg/helper"
	"gorm.io/gorm"
)

func snsShorten(db *gorm.DB) {
	var samples = []domain.SNS{
		{Url: "yt", Shorten: h.Ptr("https://www.youtube.com/"), IsPermanent: h.Ptr(true)},
		{Url: "gl", Shorten: h.Ptr("https://www.google.com/"), IsPermanent: h.Ptr(true)},
		{Url: "fb", Shorten: h.Ptr("https://www.facebook.com/"), IsPermanent: h.Ptr(false)},
		{Url: "ig", Shorten: h.Ptr("https://www.instagram.com/"), IsPermanent: h.Ptr(false)},
		{Url: "tw", Shorten: h.Ptr("https://www.twitter.com/"), IsPermanent: h.Ptr(true)},
	}
	for _, sample := range samples {
		db.Create(&sample)
	}
}
