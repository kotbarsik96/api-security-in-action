package domain

import (
	"api-security-in-action/src/models"
	"context"
)

type AuthService interface {
	ValidateSignupCredentials(login, password string) map[string]error
	CreateUser(ctx context.Context, login, password string) (*models.User, error)
	CheckCredentials(ctx context.Context, login, password string) (*models.User, error)
}

type CsrfService interface {
	GenerateToken(id string) string
	CompareToken(id string, token string) bool
}
