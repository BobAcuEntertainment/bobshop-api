package web

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	UserIDKey  = "user_id"
	RoleKey    = "role"
	IDParamKey = "id"
)

func GetUserID(c *gin.Context) uuid.UUID {
	return uuid.MustParse(c.GetString(UserIDKey))
}

func GetRole(c *gin.Context) string {
	return c.GetString(RoleKey)
}

func GetIDParam(c *gin.Context) (uuid.UUID, error) {
	param := c.Param(IDParamKey)
	id, err := uuid.Parse(param)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
