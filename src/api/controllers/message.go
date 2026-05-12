package controllers

import (
	"api-security-in-action/src/api"
	"api-security-in-action/src/api/apierrors"
	"api-security-in-action/src/domain"
	"api-security-in-action/src/models"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MessageController struct {
	Creator         domain.MessageCreator
	Repository      domain.MessageRepository
	PermissionGuard domain.PermissionGuard
}

func NewMessageController(creator domain.MessageCreator, repo domain.MessageRepository, permGuard domain.PermissionGuard) *MessageController {
	return &MessageController{
		Creator:         creator,
		Repository:      repo,
		PermissionGuard: permGuard,
	}
}

func (c *MessageController) RegisterRoutes(rootGroup *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	auditMiddleware gin.HandlerFunc) {
	auth := rootGroup.Group("", authMiddleware, auditMiddleware)
	{
		auth.POST("/spaces/:space_id/messages", c.HanldeCreateMessage)
		auth.GET("/spaces/:space_id/messages", c.HandleGetMessages)
		auth.GET("/messages/:message_id", c.HandleGetMessage)
	}
}

type MessageCreateRequest struct {
	Text string `json:"text" binding:"required"`
}

func (c *MessageController) HanldeCreateMessage(ctx *gin.Context) {
	author := ctx.MustGet("user").(models.User)

	var body MessageCreateRequest
	if err := ctx.ShouldBind(&body); err != nil {
		api.RespondError(ctx, api.Response{
			Error: api.ErrUnprocessableEntity(err.Error(), nil),
		})
		return
	}

	spaceId, err := api.GetSpaceIdParam(ctx)
	if err != nil {
		api.RespondError(ctx, api.Response{
			Error: api.ErrBadRequest(err.Error(), err),
		})
		return
	}

	hasPerm := api.CheckPermissions(ctx, c.PermissionGuard, domain.PermCreateMessageInSpace, spaceId, author.ID)
	if !hasPerm {
		return
	}

	msg, err := c.Creator.Create(ctx.Request.Context(), domain.MessageCreateData{
		SpaceID: spaceId,
		Author:  author,
		Text:    body.Text,
	})
	if err != nil {
		if errors.Is(err, apierrors.SpaceNotFound) {
			api.RespondError(ctx, api.Response{
				Error: api.ErrUnprocessableEntity(fmt.Sprintf("Space %v not found", spaceId), nil),
			})
		} else {
			api.RespondError(ctx, api.Response{
				Error: api.ErrInternal("Could not create message", err),
			})
		}

		return
	}

	api.RespondCreated(ctx, api.Response{
		Data: gin.H{"message": msg},
	})
}

type MessagesGetRequest struct {
	Since string `form:"since" time_format:"2006-01-02 15:04:05"`
}

func (c *MessageController) HandleGetMessages(ctx *gin.Context) {
	spaceId, err := api.GetSpaceIdParam(ctx)
	if err != nil {
		api.RespondError(ctx, api.Response{
			Error: api.ErrBadRequest(err.Error(), nil),
		})
		return
	}

	query := &MessagesGetRequest{}
	ctx.ShouldBindQuery(query)

	messages, err := c.Repository.GetSpaceMessages(ctx.Request.Context(), spaceId, query.Since)
	if err != nil {
		if errors.Is(err, apierrors.InvalidDatetimeString) {
			api.RespondError(ctx, api.Response{
				Error: api.ErrUnprocessableEntity(err.Error(), nil),
			})
		} else {
			api.RespondError(ctx, api.Response{
				Error: api.ErrInternal(fmt.Sprintf("Could not get messages for space id %v", spaceId), err),
			})
		}

		return
	}

	api.RespondOk(ctx, api.Response{
		Data: gin.H{"messages": messages},
	})
}

func (c *MessageController) HandleGetMessage(ctx *gin.Context) {
	id, err := api.GetMessageIdParam(ctx)
	if err != nil {
		api.RespondError(ctx, api.Response{
			Error: api.ErrBadRequest(apierrors.InvalidIdParam.Error(), nil),
		})
		return
	}

	msg, err := c.Repository.GetMessage(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.RespondError(ctx, api.Response{
				Error: api.ErrNotFound("Could not found message", nil),
			})
		} else {
			api.RespondError(ctx, api.Response{
				Error: api.ErrInternal("Could not get message", err),
			})
		}

		return
	}

	api.RespondOk(ctx, api.Response{
		Data: gin.H{"message": msg},
	})
}
