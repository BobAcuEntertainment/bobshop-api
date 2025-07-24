package web

import (
	"github.com/gin-gonic/gin"

	validator "github.com/go-playground/validator/v10"
)

func BindAndValidate(c *gin.Context, validate *validator.Validate, req any) error {
	var err error
	if c.Request.Method == "GET" {
		err = c.ShouldBindQuery(req)
	} else {
		err = c.ShouldBindJSON(req)
	}
	if err != nil {
		return err
	}
	if err := validate.Struct(req); err != nil {
		return err
	}
	return nil
}
