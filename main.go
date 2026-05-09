package main

import (
	"api-security-in-action/src/db"
	"api-security-in-action/src/domain/space"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	gormDb, err := db.NewGormDB()
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	RegisterControllers(router, gormDb)

	router.Run()
}

func RegisterControllers(router *gin.Engine, db *gorm.DB) {
	api := router.Group("/api")

	spaceCtrl := space.NewSpaceController(
		space.NewSpaceCreateService(db))

	spaceCtrl.RegisterRoutes(api)
}
