package domain

import (
	"api-security-in-action/src/models"
	"context"
)

type SpaceCreateData struct {
	Name  string
	Owner models.User
}

type SpaceCreator interface {
	Create(ctx context.Context, data SpaceCreateData) (*models.Space, error)
}
