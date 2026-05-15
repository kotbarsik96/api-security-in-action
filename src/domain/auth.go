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
	GenerateToken(sessID []byte) []byte
	CompareToken(sessID []byte, token []byte) bool
	GetCsrfProtectedMethods() []string
}
