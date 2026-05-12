package ctrlmessage

import (
	"api-security-in-action/src/api"
	"api-security-in-action/src/api/apierrors"
	"api-security-in-action/src/db/models"
	"api-security-in-action/src/domain/message"
	"context"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MessageCreator interface {
	Create(ctx context.Context, data message.MessageCreateData) (*models.Message, error)
}

type MessageRepository interface {
	GetSpaceMessages(ctx context.Context, spaceId uint, since string) ([]models.Message, error)
	GetMessage(ctx context.Context, messageId uint) (models.Message, error)
}

type MessageController struct {
	Creator    MessageCreator
	Repository MessageRepository
}

func NewMessageController(creator MessageCreator, repo MessageRepository) *MessageController {
	return &MessageController{
		Creator:    creator,
		Repository: repo,
	}
}

func (c *MessageController) RegisterRoutes(rootGroup *gin.RouterGroup, authMiddleware gin.HandlerFunc, auditMiddleware gin.HandlerFunc) {
	auth := rootGroup.Group("", authMiddleware, auditMiddleware)
	{
		auth.POST("/spaces/:space_id/messages", c.HanldeCreateMessage)
		auth.GET("/spaces/:space_id/messages", c.HandleGetMessages)
		auth.GET("/messages/:message_id", c.HandleGetMessage)
	}
}

type MessageCreateRequest struct {
	Author string `json:"author" binding:"required"`
	Text   string `json:"text" binding:"required"`
}

func (c *MessageController) HanldeCreateMessage(ctx *gin.Context) {
	var body MessageCreateRequest
	if err := ctx.ShouldBind(&body); err != nil {
		api.RespondError(ctx, api.Response{
			Error: api.ErrUnprocessableEntity(err.Error(), nil),
		})
		return
	}

	spaceId, err := GetSpaceIdParam(ctx)
	if err != nil {
		api.RespondError(ctx, api.Response{
			Error: api.ErrBadRequest(err.Error(), err),
		})
		return
	}

	msg, err := c.Creator.Create(ctx.Request.Context(), message.MessageCreateData{
		SpaceID: spaceId,
		Author:  body.Author,
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
	Since string `form:"since"`
}

func (c *MessageController) HandleGetMessages(ctx *gin.Context) {
	spaceId, err := GetSpaceIdParam(ctx)
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
	id, err := GetMessageIdParam(ctx)
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
