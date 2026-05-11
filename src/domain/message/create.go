package message

import (
	"api-security-in-action/src/api/apierrors"
	"api-security-in-action/src/db/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

type MessageCreateService struct {
	DB *gorm.DB
}

func NewMessageCreateService(db *gorm.DB) *MessageCreateService {
	return &MessageCreateService{
		DB: db,
	}
}

type MessageCreateData struct {
	SpaceID uint
	Author  string
	Text    string
}

func (s *MessageCreateService) Create(ctx context.Context, data MessageCreateData) (*models.Message, error) {
	// проверить, существует ли Space по переданному id
	_, err := gorm.G[models.Space](s.DB).Where("id = ?", data.SpaceID).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apierrors.SpaceNotFound
	}

	msg := &models.Message{
		SpaceID: data.SpaceID,
		Author:  data.Author,
		Text:    data.Text,
	}
	err = gorm.G[models.Message](s.DB).Create(ctx, msg)
	return msg, err
}
