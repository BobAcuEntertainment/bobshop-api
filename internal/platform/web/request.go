package web

import (
	"bobshop/internal/platform/response"

	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func BindAndValidate(c *gin.Context, validate *validator.Validate, req any) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, "mismatched fields", err)
		return false
	}
	if err := validate.Struct(req); err != nil {
		response.BadRequest(c, "invalid fields", err)
		return false
	}
	return true
}

func ParseUUIDFromParam(c *gin.Context, paramName string) (uuid.UUID, bool) {
	param := c.Param(paramName)
	id, err := uuid.Parse(param)
	if err != nil {
		response.BadRequest(c, "invalid id", err)
		return uuid.Nil, false
	}
	return id, true
}