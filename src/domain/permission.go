package domain

import "context"

type PermissionGuard interface {
	Check(ctx context.Context, subject EPermissionSubject, subjectID uint, userID uint) (bool, error)
}

type EPermissionSubject string

const (
	PermCreateMessageInSpace EPermissionSubject = "create_message_in_space"
	PermDeleteMessage        EPermissionSubject = "delete_message"
)

type PermissionService interface {
	Allow(ctx context.Context,
		subject EPermissionSubject,
		subjectID uint,
		userID uint) error

	Disallow(ctx context.Context,
		subject EPermissionSubject,
		subjectID uint,
		userID uint) error
}
