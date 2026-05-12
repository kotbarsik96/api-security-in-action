package message

import (
	"api-security-in-action/src/api/apierrors"
	"api-security-in-action/src/models"
	"context"
	"time"

	"gorm.io/gorm"
)

type MessagesRepository struct {
	DB *gorm.DB
}

func NewMessagesRepository(db *gorm.DB) *MessagesRepository {
	return &MessagesRepository{
		DB: db,
	}
}

func (r *MessagesRepository) GetSpaceMessages(ctx context.Context, spaceId uint, since string) ([]models.Message, error) {
	if since == "" {
		return gorm.G[models.Message](r.DB).Where("space_id = ?", spaceId).Find(ctx)
	}

	parsedSince, err := time.Parse(time.DateTime, since)
	if err != nil {
		return nil, apierrors.InvalidDatetimeString
	}

	return gorm.G[models.Message](r.DB).Where("space_id = ? AND created_at >= ?", spaceId, parsedSince).Find(ctx)
}

func (r *MessagesRepository) GetMessage(ctx context.Context, messageId uint) (models.Message, error) {
	return gorm.G[models.Message](r.DB).Where("id = ?", messageId).First(ctx)
}
