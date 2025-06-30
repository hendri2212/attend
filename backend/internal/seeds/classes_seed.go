package seeds

import (
	"attend/internal/models"
	"time"

	"gorm.io/gorm"
)

func SeedClasses(db *gorm.DB) {
	var count int64
	db.Model(&models.Class{}).Count(&count)

	if count == 0 {
		layout := "2006-01-02 15:04:05"
		t1, _ := time.Parse(layout, "2025-05-30 11:29:08")
		t2, _ := time.Parse(layout, "2025-05-30 11:28:52")

		classes := []models.Class{
			{
				ID:        1,
				Name:      "7 A",
				SchoolID:  1,
				CreatedAt: t1,
				UpdatedAt: t1,
			},
			{
				ID:        2,
				Name:      "7 B",
				SchoolID:  1,
				CreatedAt: t2,
				UpdatedAt: t2,
			},
		}

		db.Create(&classes)
	}
}
