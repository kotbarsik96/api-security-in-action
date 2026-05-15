package auth

import (
	"api-security-in-action/src/api/apierrors"
	"api-security-in-action/src/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

type PasswordService interface {
	CreateHash(password string) (string, error)
	CompareHashAndPassword(hashedPassword, password []byte) bool
}

type AuthService struct {
	DB              *gorm.DB
	PasswordService PasswordService
}

func NewAuthService(db *gorm.DB, pwdService PasswordService) *AuthService {
	return &AuthService{
		DB:              db,
		PasswordService: pwdService,
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

func (c *AuthService) CreateUser(ctx context.Context, login, password string) (*models.User, error) {
	_, err := gorm.G[models.User](c.DB).Where("login = ?", login).First(ctx)
	if err == nil {
		// логин занят
		return nil, apierrors.LoginIsTaken
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// не удалось убедиться, что логин не занят
		return nil, err
	}

	hashedPassword, err := c.PasswordService.CreateHash(password)
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

func (c *AuthService) CheckCredentials(ctx context.Context, login, password string) (*models.User, error) {
	user, err := gorm.G[models.User](c.DB).Where("login = ?", login).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierrors.InvalidCredentials
		} else {
			return nil, err
		}
	}

	passwordsMatch := c.PasswordService.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if !passwordsMatch {
		return nil, apierrors.InvalidCredentials
	}

	return &user, nil
}
