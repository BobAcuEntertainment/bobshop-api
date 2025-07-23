package web

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	UserIDKey = "user_id"
	RoleKey   = "role"
)

func GetUserID(c *gin.Context) uuid.UUID {
	return uuid.MustParse(c.GetString(UserIDKey))
}

func GetRole(c *gin.Context) string {
	return c.GetString(RoleKey)
}
