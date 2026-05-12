package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

func ErrorByCurrentEnv(abstractMessage string, err error) error {
	ce := os.Getenv("CURRENT_ENV")

	switch ce {
	case "PROD":
		return errors.New(abstractMessage)
	case "DEV":
		return fmt.Errorf("%v: %w", abstractMessage, err)
	}

	return errors.New(abstractMessage)
}

type AppError struct {
	Status int    `json:"-"`
	Code   string `json:"code"`
	// абстрактный текст сообщения об ошибке; выводится всегда
	MessageAbstract string `json:"message"`
	// текст, дополняющий MessageAbstract; выводится при DEV=1
	ErrorDescribed error `json:"-"`
}

func (err *AppError) Error() string {
	ce := os.Getenv("CURRENT_ENV")

	switch ce {
	case "PROD":
		return err.MessageAbstract
	case "DEV":
		if err.ErrorDescribed == nil {
			return fmt.Sprintf("%v (no details)", err.MessageAbstract)
		} else {
			return fmt.Sprintf("%v: %v", err.MessageAbstract, err.ErrorDescribed)
		}
	}

	return err.MessageAbstract
}

func ErrBadRequest(messageAbstract string, err error) *AppError {
	if messageAbstract == "" {
		messageAbstract = "Bad request"
	}
	return &AppError{http.StatusBadRequest, "BAD_REQUEST", messageAbstract, err}
}

func ErrUnauthorized(messageAbstract string, err error) *AppError {
	if messageAbstract == "" {
		messageAbstract = "Unauthorized"
	}
	return &AppError{http.StatusUnauthorized, "UNAUTHORIZED", messageAbstract, err}
}

func ErrForbidden(messageAbstract string, err error) *AppError {
	if messageAbstract == "" {
		messageAbstract = "Forbidden"
	}
	return &AppError{http.StatusForbidden, "FORBIDDEN", messageAbstract, err}
}

func ErrNotFound(messageAbstract string, err error) *AppError {
	if messageAbstract == "" {
		messageAbstract = "Resource not found"
	}
	return &AppError{http.StatusNotFound, "NOT_FOUND", messageAbstract, err}
}

func ErrUnprocessableEntity(messageAbstract string, err error) *AppError {
	if messageAbstract == "" {
		messageAbstract = "Unprocessable entity"
	}
	return &AppError{http.StatusUnprocessableEntity, "UNPROCESSABLE_ENTITY", messageAbstract, err}
}

func ErrInternal(messageAbstract string, err error) *AppError {
	if messageAbstract == "" {
		messageAbstract = "Internal server error"
	}
	return &AppError{http.StatusInternalServerError, "SERVER_ERROR", messageAbstract, err}
}

func ErrGatewayTimeout(messageAbstract string, err error) *AppError {
	if messageAbstract == "" {
		messageAbstract = "Gateway timeout"
	}
	return &AppError{http.StatusInternalServerError, "GATEWAY_TIMEOUT", messageAbstract, err}
}
