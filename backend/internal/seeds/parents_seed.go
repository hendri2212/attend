package seeds

import (
	"attend/internal/models"
	"time"

	"gorm.io/gorm"
)

func SeedParents(db *gorm.DB) {
	var count int64
	db.Model(&models.Parent{}).Count(&count)

	if count == 0 {
		layout := "2006-01-02 15:04:05"
		t, _ := time.Parse(layout, "2025-05-30 11:29:26")

		parent := models.Parent{
			ID:        1,
			FullName:  "Yulianoor",
			WhatsApp:  "6285746080544",
			Address:   "Kotabaru",
			CreatedAt: t,
			UpdatedAt: t,
		}

		db.Create(&parent)
	}
}
