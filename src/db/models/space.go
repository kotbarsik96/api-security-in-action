package models

import (
	"time"
)

type Space struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Owner     string    `json:"owner"`
	Messages  []Message `json:"messages"`
}
