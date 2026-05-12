package models

import "time"

type Audit struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Path      string    `json:"path"`
	UserID    uint      `json:"user_id"`
	Status    int       `json:"status"`
	UUID      string    `json:"uuid"`
}
