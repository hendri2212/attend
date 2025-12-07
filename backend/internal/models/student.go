package models

import "time"

type Student struct {
	UserID     uint      `json:"user_id" gorm:"primaryKey"`
	User       *User     `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE"`
	RFID       string    `json:"rfid" gorm:"size:50;uniqueIndex"`
	NISN       string    `json:"nisn" gorm:"size:20;uniqueIndex"`
	FullName   string    `json:"full_name" gorm:"size:100"`
	Photo      string    `json:"photo" gorm:"size:255"`
	WhatsApp   string    `json:"whatsapp" gorm:"size:20;uniqueIndex"`
	BirthPlace string    `json:"birth_place" gorm:"size:100"`
	Born       time.Time `json:"born"`
	ParentID   *uint     `json:"parent_id"`
	Parent     *Parent   `json:"parent,omitempty" gorm:"foreignKey:ParentID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	SchoolID   uint      `json:"school_id" gorm:"not null"`
	School     *School   `json:"school,omitempty" gorm:"foreignKey:SchoolID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE"`
	ClassID    uint      `json:"class_id" gorm:"not null"`
	Class      *Class    `json:"class,omitempty" gorm:"foreignKey:ClassID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type StudentResponse struct {
	ID         uint            `json:"id"`
	FullName   string          `json:"full_name"`
	WhatsApp   string          `json:"whatsapp"`
	BirthPlace string          `json:"birth_place"`
	Born       time.Time       `json:"born"`
	RFID       string          `json:"rfid"`
	NISN       string          `json:"nisn"`
	Class      *ClassResponse  `json:"class,omitempty"`
	Parent     *ParentResponse `json:"parent,omitempty"`
	Email      string          `json:"email"`
}

type ClassResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ParentResponse struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
	WhatsApp string `json:"whatsapp"`
}
