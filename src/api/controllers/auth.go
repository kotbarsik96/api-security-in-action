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

func (c *AuthController) InitAuthSession(ctx *gin.Context, userID uint) (sessions.Session, error) {
	session := sessions.Default(ctx)
	session.Set("user_id", userID)

	sessUuid, err := uuid.NewRandom()
	if err != nil {
		return session, err
	}

	session.Set("id", sessUuid.String())
	session.Save()

	return session, nil
}

func (c *AuthController) AfterSucessfullAuth(ctx *gin.Context, userID uint) error {
	session, err := c.InitAuthSession(ctx, userID)
	if err != nil {
		return err
	}

	sessID := session.Get("id").(string)
	token := c.CsrfService.GenerateToken(sessID)
	c.SetCsrfCookie(ctx, token)

	return nil
}

// ==handlers==

func (c *AuthController) RegisterRoutes(rootGroup *gin.RouterGroup, csrfMiddleware gin.HandlerFunc) {
	authGroup := rootGroup.Group("/auth")
	{
		authGroup.POST("/signup", c.HandleSignup)
		authGroup.POST("/login", c.HandleLogin)
		authGroup.GET("/csrf", c.HandleCsrfToken)
	}

	csrf := authGroup.Group("", csrfMiddleware)
	{
		csrf.POST("/logout", c.HandleLogout)
	}
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

	err = c.AfterSucessfullAuth(ctx, user.ID)

	api.RespondCreated(ctx, api.Response{
		Message: fmt.Sprintf("Welcome, %v", user.Login),
		Data: gin.H{
			"user":          user,
			"authenticated": err == nil,
		},
	})
}

func (c *AuthController) HandleCsrfToken(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userID := session.Get("user_id")
	sessID := session.Get("id")
	if userID == nil || sessID == nil {
		session.Clear()
		session.Save()
		api.RespondError(ctx, api.Response{
			Error: api.ErrUnauthorized("", nil),
		})
		return
	}

	c.InitAuthSession(ctx, userID.(uint))
	token := c.CsrfService.GenerateToken(sessID.(string))
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

	err = c.AfterSucessfullAuth(ctx, user.ID)

	if err != nil {
		api.RespondError(ctx, api.Response{
			Error: api.ErrInternal("Could not log in: server error", err),
		})
		return
	}

	api.RespondOk(ctx, api.Response{
		Message: fmt.Sprintf("Welcome, %v", user.Login),
		Data:    gin.H{"user": user},
	})
}

func (c *AuthController) HandleLogout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()

	secure := true
	if os.Getenv("CURRENT_ENV") == "DEV" {
		secure = false
	}
	ctx.SetCookie("XSRF-Token", "", -1, "/", os.Getenv("CSRF_DOMAIN"), secure, false)

	api.RespondOk(ctx, api.Response{})
}
