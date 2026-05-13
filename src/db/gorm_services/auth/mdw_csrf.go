package auth

import (
	"api-security-in-action/src/domain"
	"crypto/hmac"
	"slices"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type XCSRFTokenHeader struct {
	XCsrfToken string `json:"X-XSRF-Token"`
}

func MiddlewareCSRF(service domain.CsrfService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !slices.Contains(service.GetCsrfProtectedMethods(), strings.ToUpper(c.Request.Method)) {
			c.Next()
			return
		}

		session := sessions.Default(c)
		sessID := []byte(session.ID())

		var headers XCSRFTokenHeader
		err := c.ShouldBindHeader(&headers)
		if err != nil {
			AbortUnauthorized(c)
			return
		}

		xcsrfToken := []byte(headers.XCsrfToken)

		expected := service.GenerateToken(sessID)

		if !hmac.Equal(expected, xcsrfToken) {
			AbortUnauthorized(c)
			return
		}

		c.Next()
	}
}
