package models

import (
	"time"
)

// UserRole defines valid roles for a user.
type UserRole string

const (
	UserRoleSuperadmin UserRole = "superadmin"
	UserRoleAdmin      UserRole = "admin"
	UserRoleTeacher    UserRole = "teacher"
	UserRoleStudent    UserRole = "student"
)

type User struct {
	ID        uint     `json:"id" gorm:"primaryKey"`
	Email     string   `json:"email" gorm:"unique;size:50"`
	Password  string   `json:"password" gorm:"size:60"`
	Role      UserRole `json:"role" gorm:"type:ENUM('superadmin','admin','teacher','student');default:'student'"`
	IsActive  bool     `json:"is_active" gorm:"default:true"`
	SchoolID  uint     `json:"school_id" gorm:"not null"`
	School    *School  `json:"school,omitempty" gorm:"foreignKey:SchoolID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
