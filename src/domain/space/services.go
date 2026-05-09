package space

import (
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

func (s *SpaceCreateService) Create(data SpaceCreateData) error {
	return nil
}
