package controllers

import (
	"api-security-in-action/src/api"

	"github.com/gin-gonic/gin"
)

type CreateSpaceRequest struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

type ISpaceCreator interface {
	Create(name string, owner string) error
}

type SpaceController struct {
	Creator ISpaceCreator
}

func NewSpaceController(creator ISpaceCreator) *SpaceController {
	return &SpaceController{
		Creator: creator,
	}
}

func (sc *SpaceController) RegisterRoutes(rootGroup *gin.RouterGroup) {
	group := rootGroup.Group("/spaces")
	{
		group.POST("/", sc.Create)
	}
}

func (sc *SpaceController) Create(c *gin.Context) {
	var data CreateSpaceRequest
	err := c.ShouldBind(&data)
	if err != nil {
		api.RespondError(c, api.Response{
			Error: api.ErrUnprocessableEntity(err.Error(), nil),
		})
	}
}
