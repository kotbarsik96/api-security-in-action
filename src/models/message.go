package models

import "time"

type Message struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	SpaceID   uint      `json:"space_id"`
	AuthorID  uint      `json:"author_id"`
	Author    User      `json:"author,omitzero"`
	Text      string    `json:"text"`
}
