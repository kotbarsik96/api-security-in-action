package router

import (
	"api-security-in-action/src/api/controllers"

	"github.com/gin-gonic/gin"
)

type Router struct {
	SpaceController *controllers.SpaceController
}

func NewRouter(spaceCtrl *controllers.SpaceController) *Router {
	return &Router{
		SpaceController: spaceCtrl,
	}
}

func (r *Router) SetupRouter() *gin.Engine {
	engine := gin.Default()

	group := engine.Group("/api")

	r.SpaceController.RegisterRoutes(group)

	return engine
}
