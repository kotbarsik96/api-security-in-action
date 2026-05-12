package audit

import (
	"api-security-in-action/src/db/models"
	"context"
	"time"

	"gorm.io/gorm"
)

type AuditRepository struct {
	DB *gorm.DB
}

func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{DB: db}
}

type GetLogsFilters struct {
	Since  *time.Time
	Uuid   string
	UserID uint
	Path   string
	Status int
	Limit  int
}

func (r *AuditRepository) GetLogs(ctx context.Context, filters GetLogsFilters) ([]models.Audit, error) {
	builder := gorm.G[models.Audit](r.DB).Where("")

	if filters.Since != nil {
		builder = builder.Where("created_at >= ?", filters.Since)
	}
	if filters.Uuid != "" {
		builder = builder.Where("uuid = ?", filters.Uuid)
	}
	if filters.UserID != 0 {
		builder = builder.Where("user_id = ?", filters.UserID)
	}
	if filters.Path != "" {
		builder = builder.Where("path = ?", filters.Path)
	}
	if filters.Status != 0 {
		builder = builder.Where("status = ?", filters.Status)
	}

	limit := filters.Limit
	if limit <= 0 {
		limit = 20
	}

	return builder.Limit(limit).Find(ctx)
}
