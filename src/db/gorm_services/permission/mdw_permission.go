package permission

import (
	"api-security-in-action/src/api"
	"api-security-in-action/src/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PermissionRequest struct {
	Subject   string // операция, на которую требуется разрешение
	SubjectID uint   // конкретный индентификатор операции. Если указан 0 - не учитывается

	// требуемые права
	ReadRequired  bool
	WriteRequired bool
	DelRequired   bool
}

func AbortUnauthorized(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusForbidden, api.ErrForbidden(msg, nil))
}

func MiddlewarePermission(db *gorm.DB, request PermissionRequest) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		user := c.MustGet("user").(models.User)

		builder := gorm.G[models.Permission](db).
			Where("user_id = ?", user.ID).
			Where("subject = ?", request.Subject)
		if request.SubjectID != 0 {
			builder = builder.Where("subject_id = ?", request.SubjectID)
		}
		permission, err := builder.First(ctx)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				AbortUnauthorized(c, "Permission not found")
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrInternal("", nil))
			}

			return
		}

		if request.ReadRequired && !permission.Read {
			AbortUnauthorized(c, "Read permission required")
			return
		}

		if request.WriteRequired && !permission.Write {
			AbortUnauthorized(c, "Write permission required")
			return
		}

		if request.DelRequired && !permission.Del {
			AbortUnauthorized(c, "Del permission required")
			return
		}

		c.Next()
	}
}
