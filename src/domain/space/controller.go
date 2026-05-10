package space

import (
	"api-security-in-action/src/api"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

type SpaceCreator interface {
	Create(ctx context.Context, data SpaceCreateData) error
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

	err = c.Creator.Create(ctx.Request.Context(), data)

	if err != nil {
		api.RespondError(ctx, api.Response{
			Error: api.ErrInternal("Could not create space", err),
		})
		return
	}

	api.RespondCreated(ctx, api.Response{
		Message: fmt.Sprintf("Created new space \"%v\"", data.Name),
	})
}
