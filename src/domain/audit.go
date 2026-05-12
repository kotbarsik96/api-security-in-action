package domain

import (
	"api-security-in-action/src/models"
	"context"
	"time"
)

type GetAuditLogsFilters struct {
	Since  *time.Time
	Uuid   string
	UserID uint
	Path   string
	Status int
	Limit  int
}

type AuditRepository interface {
	GetLogs(ctx context.Context, filters GetAuditLogsFilters) ([]models.Audit, error)
}
