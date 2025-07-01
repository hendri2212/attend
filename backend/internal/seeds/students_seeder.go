package seeds

import (
	"attend/internal/models"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func uintPtr(i uint) *uint {
	return &i
}

func SeedStudents(db *gorm.DB) {
	var count int64
	db.Model(&models.Student{}).Count(&count)
	if count == 0 {
		seeds := []struct {
			Email      string
			Password   string
			NISN       string
			FullName   string
			Photo      string
			WhatsApp   string
			BirthPlace string
			Born       time.Time
			ParentID   *uint
			SchoolID   uint
			ClassID    uint
			RFID       string
		}{
			{
				Email:      "john.doe@example.com",
				Password:   "password123",
				NISN:       "1234567890",
				FullName:   "John Doe",
				Photo:      "johndoe.jpg",
				WhatsApp:   "081234567890",
				BirthPlace: "City A",
				Born:       time.Date(2005, time.January, 15, 0, 0, 0, 0, time.UTC),
				ParentID:   uintPtr(1),
				SchoolID:   1,
				ClassID:    1,
				RFID:       "RFID123456",
			},
			{
				Email:      "jane.smith@example.com",
				Password:   "password123",
				NISN:       "0987654321",
				FullName:   "Jane Smith",
				Photo:      "janesmith.jpg",
				WhatsApp:   "089876543210",
				BirthPlace: "City B",
				Born:       time.Date(2006, time.February, 20, 0, 0, 0, 0, time.UTC),
				ParentID:   uintPtr(1),
				SchoolID:   1,
				ClassID:    1,
				RFID:       "RFID654321",
			},
		}

		for _, s := range seeds {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
			user := models.User{
				Email:    s.Email,
				Password: string(hashedPassword),
				Role:     "student",
				SchoolID: s.SchoolID,
			}
			db.Create(&user)

			student := models.Student{
				UserID:     user.ID,
				NISN:       s.NISN,
				FullName:   s.FullName,
				Photo:      s.Photo,
				WhatsApp:   s.WhatsApp,
				BirthPlace: s.BirthPlace,
				Born:       s.Born,
				ParentID:   s.ParentID,
				ClassID:    s.ClassID,
				RFID:       s.RFID,
				SchoolID:   s.SchoolID,
			}
			db.Create(&student)
		}
	}
}
