package db

import (
	"api-security-in-action/src/domain/message"
	"api-security-in-action/src/domain/space"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&message.Message{}, &space.Space{})

	return db, nil
}
