package space

import (
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

func (s *SpaceCreateService) Create(ctx context.Context, data SpaceCreateData) error {
	space := &Space{
		Name:  data.Name,
		Owner: data.Owner,
	}

	err := gorm.G[Space](s.DB).Create(ctx, space)

	return err
}
