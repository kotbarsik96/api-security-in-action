package db

import (
	"api-security-in-action/src/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewGormDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Space{}, &models.Message{}, &models.User{}, &models.Audit{})

	return db, nil
}
