package space

import (
	"api-security-in-action/src/domain"
	"api-security-in-action/src/models"
	"context"

	"gorm.io/gorm"
)

type SpaceCreateService struct {
	DB *gorm.DB
}

func NewSpaceCreateService(db *gorm.DB) *SpaceCreateService {
	return &SpaceCreateService{
		DB: db,
	}
}

func (s *SpaceCreateService) Create(ctx context.Context, data domain.SpaceCreateData) (*models.Space, error) {
	space := &models.Space{
		Name:    data.Name,
		OwnerID: data.Owner.ID,
	}

	err := gorm.G[models.Space](s.DB).Create(ctx, space)

	return space, err
}
