package ctrlaudit

import (
	"api-security-in-action/src/api"
	"api-security-in-action/src/db/models"
	"api-security-in-action/src/domain/audit"
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

type AuditRepository interface {
	GetLogs(ctx context.Context, filters audit.GetLogsFilters) ([]models.Audit, error)
}

type AuditController struct {
	Repository AuditRepository
}

func NewAuditController(repo AuditRepository) *AuditController {
	return &AuditController{
		Repository: repo,
	}
}

func (c *AuditController) RegisterRoutes(rootGroup *gin.RouterGroup, authMiddleware gin.HandlerFunc, auditMiddleware gin.HandlerFunc) {
	auth := rootGroup.Group("", authMiddleware, auditMiddleware)
	{
		auth.GET("/audit/logs", c.HandleGetLogs)
	}
}

type LogsFilters struct {
	Since  *time.Time `form:"since" time_format:"2006-01-02 15:04:05"`
	Uuid   string     `form:"uuid"`
	UserID uint       `form:"user_id"`
	Path   string     `form:"path"`
	Status int        `form:"status"`
	Limit  int        `form:"limit"`
}

func (c *AuditController) HandleGetLogs(ctx *gin.Context) {
	query := LogsFilters{}
	err := ctx.ShouldBindQuery(&query)
	if err != nil {
		api.RespondError(ctx, api.Response{
			Error: api.ErrUnprocessableEntity(err.Error(), nil),
		})
		return
	}

	logs, err := c.Repository.GetLogs(ctx.Request.Context(), audit.GetLogsFilters{
		Since:  query.Since,
		Uuid:   query.Uuid,
		UserID: query.UserID,
		Path:   query.Path,
		Status: query.Status,
		Limit:  query.Limit,
	})

	if err != nil {
		api.RespondError(ctx, api.Response{
			Error: api.ErrInternal("Could not get logs", err),
		})
		return
	}

	api.RespondOk(ctx, api.Response{
		Data: gin.H{"logs": logs},
	})
}
