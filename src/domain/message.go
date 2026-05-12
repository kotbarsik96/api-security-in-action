package domain

import (
	"api-security-in-action/src/models"
	"context"
)

type MessageCreateData struct {
	SpaceID uint
	Author  models.User
	Text    string
}

type MessageCreator interface {
	Create(ctx context.Context, data MessageCreateData) (*models.Message, error)
}

type MessageRepository interface {
	GetSpaceMessages(ctx context.Context, spaceId uint, since string) ([]models.Message, error)
	GetMessage(ctx context.Context, messageId uint) (models.Message, error)
}
