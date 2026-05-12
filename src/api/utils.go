package api

import (
	"api-security-in-action/src/api/apierrors"
	"api-security-in-action/src/domain"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetIdParam(c *gin.Context, paramName string) (uint, error) {
	id, err := strconv.ParseUint(c.Param(paramName), 10, 0)
	if err != nil {
		return 0, apierrors.InvalidIdParam
	}

	return uint(id), nil
}

func GetSpaceIdParam(c *gin.Context) (uint, error) {
	return GetIdParam(c, "space_id")
}

func GetMessageIdParam(c *gin.Context) (uint, error) {
	return GetIdParam(c, "message_id")
}

// проверяет права и автоматически записывает ошибку в gin.Context, если она возникла. Клиенту нужно проверить только возвращаемое значение
func CheckPermissions(c *gin.Context, guard domain.PermissionGuard, subject domain.EPermissionSubject, subjectID uint, userID uint) bool {
	ctx := c.Request.Context()
	hasPerm, err := guard.Check(ctx, subject, subjectID, userID)

	if err != nil {
		RespondError(c, Response{
			Error: ErrInternal("", err),
		})
		return false
	}

	if !hasPerm {
		RespondForbidden(c)
	}
	return hasPerm
}
