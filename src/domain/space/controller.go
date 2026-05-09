package space

import (
	"api-security-in-action/src/api"

	"github.com/gin-gonic/gin"
)

type SpaceCreator interface {
	Create(SpaceCreateData) error
}

type SpaceCreateData struct {
	Name  string `json:"name" binding:"required"`
	Owner string `json:"owner" binding:"required"`
}

type SpaceController struct {
	Creator SpaceCreator
}

func NewSpaceController(creator SpaceCreator) *SpaceController {
	return &SpaceController{
		Creator: creator,
	}
}

func (c *SpaceController) RegisterRoutes(rootRoute *gin.RouterGroup) {
	rootRoute.POST("/spaces", c.HandleCreateSpace)
}

func (c *SpaceController) HandleCreateSpace(ctx *gin.Context) {
	var data SpaceCreateData
	err := ctx.ShouldBind(&data)
	if err != nil {
		api.RespondError(ctx, api.Response{
			Error: api.ErrUnprocessableEntity(err.Error(), nil),
		})
		return
	}

	c.Creator.Create(data)
}
