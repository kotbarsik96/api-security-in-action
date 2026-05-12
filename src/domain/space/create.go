package space

import (
	"api-security-in-action/src/db/models"
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

type SpaceCreateData struct {
	Name  string
	Owner models.User
}

func (s *SpaceCreateService) Create(ctx context.Context, data SpaceCreateData) (*models.Space, error) {
	space := &models.Space{
		Name:    data.Name,
		OwnerID: data.Owner.ID,
	}

	err := gorm.G[models.Space](s.DB).Create(ctx, space)

	return space, err
}
