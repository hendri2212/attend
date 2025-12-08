package models

import (
	"time"

	"gorm.io/gorm"
)

type Class struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:50;unique"`
	SchoolID  uint      `json:"school_id" gorm:"not null"`
	School    School    `json:"school" gorm:"foreignKey:SchoolID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeDelete is a GORM hook that runs before a Class is deleted.
// It deletes all Users associated with students in this class.
// This is necessary because GORM's cascading delete for Students might run
// via DB constraints (which don't fire hooks) or via GORM (which might not fire Student hooks depending on how it's done).
// Safest is to explicitly cleanup Users of students in this class.
func (c *Class) BeforeDelete(tx *gorm.DB) (err error) {
	// Find all students in this class
	var students []Student
	if err := tx.Where("class_id = ?", c.ID).Find(&students).Error; err != nil {
		return err
	}

	// Delete users for these students
	for _, student := range students {
		if student.UserID != 0 {
			if err := tx.Delete(&User{}, student.UserID).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
