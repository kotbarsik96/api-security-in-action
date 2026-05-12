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

type SpaceController struct {
	Creator SpaceCreator
}

func NewSpaceController(creator SpaceCreator) *SpaceController {
	return &SpaceController{
		Creator: creator,
	}
}

func (c *SpaceController) RegisterRoutes(rootRoute *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	auditMiddleware gin.HandlerFunc) {
	auth := rootRoute.Group("", authMiddleware, auditMiddleware)
	{
		auth.POST("/spaces", c.HandleCreateSpace)
	}
}

type SpaceCreateRequest struct {
	Name string `json:"name" binding:"required"`
}

func (c *SpaceController) HandleCreateSpace(ctx *gin.Context) {
	user := ctx.MustGet("user").(models.User)

	var body SpaceCreateRequest
	err := ctx.ShouldBind(&body)
	if err != nil {
		api.RespondError(ctx, api.Response{
			Error: api.ErrUnprocessableEntity(err.Error(), nil),
		})
		return
	}

	s, err := c.Creator.Create(ctx.Request.Context(), space.SpaceCreateData{
		Name:  body.Name,
		Owner: user,
	})

	if err != nil {
		api.RespondError(ctx, api.Response{
			Error: api.ErrInternal("Could not create space", err),
		})
		return
	}

	api.RespondCreated(ctx, api.Response{
		Message: fmt.Sprintf("Created new space \"%v\"", body.Name),
		Data:    gin.H{"space": s},
	})
}
