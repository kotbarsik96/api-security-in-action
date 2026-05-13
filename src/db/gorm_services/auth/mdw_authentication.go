package auth

import (
	"api-security-in-action/src/api"
	"api-security-in-action/src/models"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHeader struct {
	Authorization string
}

func AbortUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrUnauthorized("", nil))
}

func MiddlewareAuthentication(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userId := session.Get("user_id")
		if userId == nil {
			AbortUnauthorized(c)
			return
		}

		user, err := gorm.G[models.User](db).Where("id = ?", userId).First(c.Request.Context())
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				session.Clear()
				AbortUnauthorized(c)
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrInternal("", nil))
			}

			return
		}

		c.Set("user", user)

		c.Next()
	}
}

func MiddlewareAuthenticationDumb(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// временная реализация до прочтения главы 4 (Session Cookie Authentication)
		var h AuthHeader
		if err := ctx.ShouldBindHeader(&h); err != nil {
			AbortUnauthorized(ctx)
			return
		}

		prefix := "Dumb "
		length := len(prefix)

		if len(h.Authorization) <= length || !strings.HasPrefix(h.Authorization, prefix) {
			AbortUnauthorized(ctx)
			return
		}

		val := h.Authorization[len(prefix):]
		creds := strings.Split(val, "|=|")
		if len(creds) != 2 {
			AbortUnauthorized(ctx)
			return
		}

		login := creds[0]
		password := creds[1]

		user, err := gorm.G[models.User](db).Where("login = ?", login).First(ctx.Request.Context())
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				AbortUnauthorized(ctx)
			} else {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrInternal("", nil))
			}

			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			AbortUnauthorized(ctx)
			return
		}

		ctx.Set("user", user)

		ctx.Next()
	}
}
