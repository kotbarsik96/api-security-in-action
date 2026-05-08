package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Data    any       `json:"data,omitempty"`
	Message string    `json:"message,omitempty"`
	Error   *AppError `json:"error,omitempty"`
}

func RespondOkWithCode(c *gin.Context, response Response, statusCode int) {
	r := gin.H{
		"ok": true,
	}
	if response.Data != nil {
		r["data"] = response.Data
	}
	if response.Message != "" {
		r["message"] = response.Message
	}
	c.JSON(statusCode, r)
}

func RespondOk(c *gin.Context, response Response) {
	RespondOkWithCode(c, response, http.StatusOK)
}

func RespondCreated(c *gin.Context, response Response) {
	RespondOkWithCode(c, response, http.StatusCreated)
}

func RespondError(c *gin.Context, response Response) {
	c.JSON(response.Error.Status, gin.H{
		"ok":    false,
		"error": response.Error.Error(),
	})
}
