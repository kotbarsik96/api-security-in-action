package permission

import (
	"api-security-in-action/src/domain"
	"api-security-in-action/src/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

type PermissionGuard struct {
	DB *gorm.DB
}

func NewPermissionGuard(db *gorm.DB) *PermissionGuard {
	return &PermissionGuard{
		DB: db,
	}
}

func (g *PermissionGuard) Check(ctx context.Context,
	subject domain.EPermissionSubject,
	subjectID uint,
	userID uint) (bool, error) {

	builder := gorm.G[models.Permission](g.DB).
		Where("subject = ?", subject).
		Where("user_id = ?", userID)
	if subjectID != 0 {
		builder = builder.Where("subject_id = ?", subjectID)
	}

	p, err := builder.First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return p.Allowed, nil
}
