package auth

import (
	"api-security-in-action/src/api/apierrors"
	"api-security-in-action/src/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		DB: db,
	}
}

func (c *AuthService) ValidateSignupCredentials(login, password string) map[string]error {
	errs := make(map[string]error)

	if err := ValidateLogin(login); err != nil {
		errs["login"] = err
	}

	if err := ValidatePassword(password); err != nil {
		errs["password"] = err
	}

	return errs
}

func (c *AuthService) Signup(ctx context.Context, login, password string) (*models.User, error) {
	_, err := gorm.G[models.User](c.DB).Where("login = ?", login).First(ctx)
	exists := !errors.Is(err, gorm.ErrRecordNotFound)
	if exists {
		return nil, apierrors.LoginIsTaken
	}

	if err != nil {
		return nil, err
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Login:    login,
		Password: hashedPassword,
	}
	err = gorm.G[models.User](c.DB).Create(ctx, user)

	return user, nil
}

func (c *AuthService) Login(ctx context.Context, login, password string) (*models.User, error) {
	return nil, nil
}
