package models

import "time"

type Parent struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	FullName  string `json:"full_name" gorm:"size:100"`
	WhatsApp  string `json:"whatsapp" gorm:"size:20;uniqueIndex"`
	Address   string `json:"address" gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
