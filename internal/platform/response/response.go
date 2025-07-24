package response

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Data    any       `json:"data,omitempty"`
	Error   *APIError `json:"error,omitempty"`
}

type APIError struct {
	Code   string `json:"code,omitempty"`
	Detail any    `json:"detail,omitempty"`
}

var showErrorDetails bool

func Configure(showDetails bool) {
	showErrorDetails = showDetails
}

func Success(c *gin.Context, statusCode int, message string, data any) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SimpleSuccess(c *gin.Context, message string) {
	Success(c, http.StatusOK, message, nil)
}

func Created(c *gin.Context, message string, data any) {
	Success(c, http.StatusCreated, message, data)
}

func NoContent(c *gin.Context, message string) {
	Success(c, http.StatusNoContent, message, nil)
}

func Error(c *gin.Context, statusCode int, errorCode string, clientMessage string, err error) {
	log.Printf("ERROR [%s]: %v", errorCode, err) // log for debugging

	var detail any
	if showErrorDetails && err != nil {
		detail = err.Error()
	}

	c.AbortWithStatusJSON(statusCode, Response{
		Success: false,
		Message: clientMessage,
		Error: &APIError{
			Code:   errorCode,
			Detail: detail,
		},
	})
}

func BadRequest(c *gin.Context, detail string, err error) {
	Error(c, http.StatusBadRequest, "BAD_REQUEST", detail, err)
}

func Unauthorized(c *gin.Context, err error) {
	Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized", err)
}

func Forbidden(c *gin.Context, err error) {
	Error(c, http.StatusForbidden, "FORBIDDEN", "Forbidden", err)
}

func NotFound(c *gin.Context, err error) {
	Error(c, http.StatusNotFound, "NOT_FOUND", "Not found", err)
}

func Conflict(c *gin.Context, detail string, err error) {
	Error(c, http.StatusConflict, "CONFLICT", detail, err)
}

func InternalError(c *gin.Context, err error) {
	Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Internal server error", err)
}
