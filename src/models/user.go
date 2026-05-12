package models

import "time"

type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Login     string    `json:"login"`
	Password  string    `json:"-"`
}
