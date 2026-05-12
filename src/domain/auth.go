package domain

import (
	"api-security-in-action/src/models"
	"context"
)

type AuthService interface {
	ValidateSignupCredentials(login, password string) map[string]error
	Signup(ctx context.Context, login, password string) (*models.User, error)
	Login(ctx context.Context, login, password string) (*models.User, error)
}
