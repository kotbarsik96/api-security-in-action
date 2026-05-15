package controllers

import (
	"api-security-in-action/src/api"
	"api-security-in-action/src/api/apierrors"
	"api-security-in-action/src/domain"
	"errors"
	"fmt"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ==controller==

type AuthController struct {
	AuthService domain.AuthService
	CsrfService domain.CsrfService
}

func NewAuthController(authService domain.AuthService, csrfService domain.CsrfService) *AuthController {
	return &AuthController{
		AuthService: authService,
		CsrfService: csrfService,
	}
}

// ==helpers==

func (c *AuthController) SetCsrfCookie(ctx *gin.Context, token string) {
	secure := true
	if os.Getenv("CURRENT_ENV") == "DEV" {
		secure = false
	}

	maxAge := 3600 * 24

	ctx.SetCookie("XSRF-Token", token, maxAge, "/", os.Getenv("CSRF_DOMAIN"), secure, false)
}

func (c *AuthController) InitAuthSession(ctx *gin.Context, userID uint) sessions.Session {
	session := sessions.Default(ctx)
	session.Set("user_id", userID)
	session.Set("id", uuid.New().String())
	session.Save()
	return session
}

func (c *AuthController) AfterSucessfullAuth(ctx *gin.Context, userID uint) {
	session := c.InitAuthSession(ctx, userID)
	sessID := session.Get("id").(string)
	token := c.CsrfService.GenerateToken(sessID)
	c.SetCsrfCookie(ctx, token)
}

// ==handlers==

func (c *AuthController) RegisterRoutes(rootGroup *gin.RouterGroup) {
	rootGroup.POST("/auth/signup", c.HandleSignup)
	rootGroup.POST("/auth/login", c.HandleLogin)
	rootGroup.GET("/csrf", c.HandleCsrfToken)
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
			Error: api.ErrUnprocessableEntity("Could not sign up: invalid credentials", nil),
			Data:  gin.H{"errors": errs},
		})
		return
	}

	user, err := c.AuthService.CreateUser(ctx.Request.Context(), data.Login, data.Password)
	if err != nil {
		if errors.Is(err, apierrors.LoginIsTaken) {
			api.RespondError(ctx, api.Response{
				Error: api.ErrUnprocessableEntity(err.Error(), nil),
			})
		} else {
			api.RespondError(ctx, api.Response{
				Error: api.ErrInternal("Could not sign up", err),
			})
		}

		return
	}

	c.AfterSucessfullAuth(ctx, user.ID)

	api.RespondCreated(ctx, api.Response{
		Message: fmt.Sprintf("Welcome, %v", user.Login),
		Data:    gin.H{"user": user},
	})
}

func (c *AuthController) HandleCsrfToken(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userID := session.Get("user_id")
	if userID == nil {
		api.RespondError(ctx, api.Response{
			Error: api.ErrUnauthorized("", nil),
		})
		return
	}

	c.InitAuthSession(ctx, userID.(uint))
	token := c.CsrfService.GenerateToken(session.Get("id").(string))
	c.SetCsrfCookie(ctx, token)

	api.RespondOk(ctx, api.Response{})
}

type UserLoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (c *AuthController) HandleLogin(ctx *gin.Context) {
	var body UserLoginRequest
	err := ctx.ShouldBind(&body)
	if err != nil {
		api.RespondError(ctx, api.Response{
			Error: api.ErrUnprocessableEntity(err.Error(), nil),
		})
		return
	}

	user, err := c.AuthService.CheckCredentials(ctx.Request.Context(), body.Login, body.Password)
	if err != nil {
		if errors.Is(err, apierrors.InvalidCredentials) {
			api.RespondError(ctx, api.Response{
				Error: api.ErrUnauthorized(err.Error(), nil),
			})
		} else {
			api.RespondError(ctx, api.Response{
				Error: api.ErrInternal("Could not login", err),
			})
		}

		return
	}

	c.AfterSucessfullAuth(ctx, user.ID)

	api.RespondOk(ctx, api.Response{
		Message: fmt.Sprintf("Welcome, %v", user.Login),
		Data:    gin.H{"user": user},
	})
}
