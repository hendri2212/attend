package models

import "time"

type Teacher struct {
	UserID   uint    `json:"user_id" gorm:"primaryKey"`
	User     *User   `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE"`
	FullName string  `json:"full_name" gorm:"size:50"`
	Nip      *string `json:"nip" gorm:"size:25;default:NULL"`
	// PositionID *uint     `json:"position_id" gorm:"default:NULL"`
	// Position   *Position `json:"position" gorm:"foreignKey:PositionID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	Photo     *string   `json:"photo" gorm:"default:NULL"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
