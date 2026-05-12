package audit

import (
	"api-security-in-action/src/db/models"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AuditMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		user := c.MustGet("user").(models.User)
		uuid, err := uuid.NewUUID()
		auditUuid := uuid.String()

		err = gorm.G[models.Audit](db).Create(ctx, &models.Audit{
			Path:   c.Request.URL.Path,
			UserID: user.ID,
			UUID:   auditUuid,
		})

		if err != nil {
			log.Printf("Error while saving audit log for request %v start: %v\n", auditUuid, err)
		}

		c.Next()

		err = gorm.G[models.Audit](db).Create(ctx, &models.Audit{
			Path:   c.Request.URL.Path,
			UserID: user.ID,
			Status: c.Writer.Status(),
			UUID:   auditUuid,
		})

		if err != nil {
			log.Printf("Error while saving audit log for request %v end: %v\n", auditUuid, err)
		}
	}
}
