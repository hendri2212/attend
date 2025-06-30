package models

import "time"

type Class struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:50;unique"`
	SchoolID  uint      `json:"school_id" gorm:"not null"`
	School    School    `json:"school" gorm:"foreignKey:SchoolID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
