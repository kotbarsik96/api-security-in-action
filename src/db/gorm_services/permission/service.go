package permission

import (
	"api-security-in-action/src/domain"
	"api-security-in-action/src/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

type PermissionService struct {
	DB *gorm.DB
}

func NewPermissionService(db *gorm.DB) *PermissionService {
	return &PermissionService{
		DB: db,
	}
}

func (s *PermissionService) Update(ctx context.Context,
	subject domain.EPermissionSubject,
	subjectID uint,
	userID uint,
	allowed bool) error {
	builder := gorm.G[models.Permission](s.DB).
		Where("user_id = ?", userID).
		Where("subject = ?", subject)

	if subjectID != 0 {
		builder = builder.Where("subject_id = ?", subjectID)
	}

	p, err := builder.First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if !allowed {
			return nil
		}

		return s.allowFirstTime(ctx, subject, subjectID, userID)
	} else if err != nil {
		return err
	}

	if p.Allowed == allowed {
		return nil
	}

	_, err = builder.Update(ctx, "allowed", allowed)
	return err
}

func (s *PermissionService) allowFirstTime(ctx context.Context, subject domain.EPermissionSubject, subjectID uint, userID uint) error {
	perm := &models.Permission{
		Subject: string(subject),
		UserID:  userID,
		Allowed: true,
	}
	if subjectID != 0 {
		perm.SubjectID = subjectID
	}

	return gorm.G[models.Permission](s.DB).Create(ctx, perm)
}

func (s *PermissionService) Allow(ctx context.Context, subject domain.EPermissionSubject, subjectID uint, userID uint) error {
	return s.Update(ctx, subject, subjectID, userID, true)
}

func (s *PermissionService) Disallow(ctx context.Context, subject domain.EPermissionSubject, subjectID uint, userID uint) error {
	return s.Update(ctx, subject, subjectID, userID, false)
}
