package middleware

import (
	"github.com/gin-gonic/gin"

	"bobshop/internal/platform/response"
	"bobshop/internal/platform/security"
)

const (
	accessTokenCookieName = "access_token"
)

func AuthMiddleware(parser security.Tokenizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := c.Cookie(accessTokenCookieName)
		claims, err := parser.ParseToken(token)
		if err != nil {
			response.Unauthorized(c, err)
			return
		}
		c.Set("userID", claims["sub"])
		c.Set("role", claims["role"])
	}
}
