package space

import (
	"gorm.io/gorm"
)

type GormSpaceCreator struct {
	db *gorm.DB
}

func NewGormSpaceCreator(db *gorm.DB) *GormSpaceCreator {
	return &GormSpaceCreator{
		db: db,
	}
}

func (c *GormSpaceCreator) Create(name string, owner string) error {

	return nil
}
