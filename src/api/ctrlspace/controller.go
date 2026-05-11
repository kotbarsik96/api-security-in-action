package ctrlspace

import (
	"api-security-in-action/src/api"
	"api-security-in-action/src/db/models"
	"api-security-in-action/src/domain/space"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

type SpaceCreator interface {
	Create(ctx context.Context, data space.SpaceCreateData) (*models.Space, error)
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

	s, err := c.Creator.Create(ctx.Request.Context(), space.SpaceCreateData{
		Name:  data.Name,
		Owner: data.Owner,
	})

	if err != nil {
		api.RespondError(ctx, api.Response{
			Error: api.ErrInternal("Could not create space", err),
		})
		return
	}

	api.RespondCreated(ctx, api.Response{
		Message: fmt.Sprintf("Created new space \"%v\"", data.Name),
		Data:    gin.H{"space": s},
	})
}
