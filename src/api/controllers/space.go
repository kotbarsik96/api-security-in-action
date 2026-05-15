package controllers

import (
	"api-security-in-action/src/api"
	"api-security-in-action/src/domain"
	"api-security-in-action/src/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

type SpaceController struct {
	Creator           domain.SpaceCreator
	PermissionService domain.PermissionService
}

func NewSpaceController(creator domain.SpaceCreator, permService domain.PermissionService) *SpaceController {
	return &SpaceController{
		Creator:           creator,
		PermissionService: permService,
	}
}

func (c *SpaceController) RegisterRoutes(rootRoute *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	auditMiddleware gin.HandlerFunc,
	csrfMiddleware gin.HandlerFunc) {
	auth := rootRoute.Group("", authMiddleware, auditMiddleware)
	csrf := auth.Group("", csrfMiddleware)
	{
		csrf.POST("/spaces", c.HandleCreateSpace)
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

	s, err := c.Creator.Create(ctx.Request.Context(), domain.SpaceCreateData{
		Name:  body.Name,
		Owner: user,
	})

	if err != nil {
		api.RespondError(ctx, api.Response{
			Error: api.ErrInternal("Could not create space", err),
		})
		return
	}

	c.PermissionService.Allow(ctx.Request.Context(), domain.PermCreateMessageInSpace, s.ID, user.ID)

	api.RespondCreated(ctx, api.Response{
		Message: fmt.Sprintf("Created new space \"%v\"", body.Name),
		Data:    gin.H{"space": s},
	})
}
