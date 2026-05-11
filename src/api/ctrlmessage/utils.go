package ctrlmessage

import (
	"api-security-in-action/src/api/apierrors"
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
