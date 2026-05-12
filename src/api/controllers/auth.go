package controllers

import (
	"api-security-in-action/src/api"
	"api-security-in-action/src/domain"
	"fmt"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthService domain.AuthService
}

func NewAuthController(authService domain.AuthService) *AuthController {
	return &AuthController{
		AuthService: authService,
	}
}

func (c *AuthController) RegisterRoutes(rootGroup *gin.RouterGroup) {
	rootGroup.POST("/auth/signup", c.HandleSignup)
}

type UserSignupRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (c *AuthController) HandleSignup(ctx *gin.Context) {
	var data UserSignupRequest
	if err := ctx.ShouldBind(&data); err != nil {
		api.RespondError(ctx, api.Response{
			Error: api.ErrUnprocessableEntity(err.Error(), nil),
		})
		return
	}

	errs := c.AuthService.ValidateSignupCredentials(data.Login, data.Password)
	if len(errs) > 0 {
		api.RespondError(ctx, api.Response{
			Error: api.ErrUnprocessableEntity("Could not sign up: validation errors", nil),
			Data:  gin.H{"errors": errs},
		})
		return
	}

	user, err := c.AuthService.Signup(ctx.Request.Context(), data.Login, data.Password)
	if err != nil {
		api.RespondError(ctx, api.Response{
			Error: api.ErrInternal("Could not sign up", err),
		})
		return
	}

	api.RespondCreated(ctx, api.Response{
		Message: fmt.Sprintf("Welcome, %v", user.Login),
		Data:    gin.H{"user": user},
	})
}
