package models

import "time"

type Attendance struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	StudentID  uint       `json:"student_id" gorm:"not null;index:uniq_attendance_day,unique"`
	Student    *Student   `json:"student,omitempty" gorm:"foreignKey:StudentID;references:UserID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE"`
	Date       time.Time  `json:"date" gorm:"type:date;not null;index:uniq_attendance_day,unique"` // calendar day of the attendance
	CheckInAt  *time.Time `json:"check_in_at,omitempty"`
	CheckOutAt *time.Time `json:"check_out_at,omitempty"`
	Method     string     `json:"method" gorm:"size:50"` // e.g. rfid, manual
	Note       string     `json:"note" gorm:"size:255"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
