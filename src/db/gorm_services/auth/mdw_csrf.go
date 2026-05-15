package auth

import (
	"api-security-in-action/src/domain"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func MiddlewareCSRF(service domain.CsrfService) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessID := session.Get("id").(string)

		xcsrfToken := c.GetHeader("X-XSRF-Token")
		if xcsrfToken == "" {
			AbortUnauthorized(c)
			return
		}

		if !service.CompareToken(sessID, xcsrfToken) {
			AbortUnauthorized(c)
			return
		}

		c.Next()
	}
}
