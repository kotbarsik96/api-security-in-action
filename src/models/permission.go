package models

import "time"

type Permission struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Subject   string    `json:"subject"`
	SubjectID uint      `json:"subject_id"`
	UserID    uint      `json:"user_id"`
	User      User      `json:"user,omitzero"`
	Allowed   bool      `json:"allowed"`
}
