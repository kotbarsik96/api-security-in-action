package models

import (
	"time"
)

type Space struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	OwnerID   uint      `json:"owner_id"`
	Owner     User      `json:"owner,omitzero"`
	Messages  []Message `json:"messages"`
}
