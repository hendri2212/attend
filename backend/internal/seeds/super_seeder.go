package seeds

import (
	"time"

	"attend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedSuperAdmin membuat user superadmin dengan nip, position_id, leader_id null
func SeedSuperAdmin(db *gorm.DB) error {
	// Cek dulu apakah superadmin sudah ada
	var count int64
	db.Model(&models.User{}).
		Where("role = ?", "superadmin").
		Count(&count)
	if count > 0 {
		return nil // sudah pernah di‐seed
	}

	// Hash default password (ganti dengan env var jika perlu)
	pwd, err := bcrypt.GenerateFromPassword([]byte("Secret@123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	super := models.User{
		Email:     "arifin.hendri465@gmail.com",
		Password:  string(pwd),
		Role:      models.UserRoleSuperadmin,
		SchoolID:  1, // Asumsi ada school dengan ID 1
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return db.Create(&super).Error
}
