package models

import "time"

type School struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	NPSN      string `json:"npsn" gorm:"size:50;unique"`
	Name      string `json:"name" gorm:"size:200;unique"`
	Address   string `json:"address" gorm:"size:255"`
	Logo      string `json:"logo" gorm:"size:255"`
	Email     string `json:"email" gorm:"size:100"`
	Telp      string `json:"telp" gorm:"size:20"`
	Message   string `json:"message" gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
