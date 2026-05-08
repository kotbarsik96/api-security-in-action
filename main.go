package main

import (
	"api-security-in-action/src/api/controllers"
	"api-security-in-action/src/db"
	"api-security-in-action/src/domain/space"
	"api-security-in-action/src/router"
)

func main() {
	sdb, err := db.SetupDatabase()
	if err != nil {
		panic(err)
	}

	spaceController := controllers.NewSpaceController(space.NewGormSpaceCreator(sdb))

	router := router.NewRouter(spaceController)
	engine := router.SetupRouter()
	engine.Run()
}
