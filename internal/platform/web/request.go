package web

import (
	"github.com/gin-gonic/gin"

	validator "github.com/go-playground/validator/v10"

	"bobshop/internal/platform/response"
)

func BindAndValidate(c *gin.Context, validate *validator.Validate, req any) bool {
	var err error
	if c.Request.Method == "GET" {
		err = c.ShouldBindQuery(req)
	} else {
		err = c.ShouldBindJSON(req)
	}
	if err != nil {
		response.BadRequest(c, "mismatched fields", err)
		return false
	}
	if err := validate.Struct(req); err != nil {
		response.BadRequest(c, "invalid fields", err)
		return false
	}
	return true
}