package message

import (
	"time"
)

type Message struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	SpaceID   uint      `json:"space_id"`
}
