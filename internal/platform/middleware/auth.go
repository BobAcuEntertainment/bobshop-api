package middleware

import (
	"github.com/gin-gonic/gin"

	"bobshop/internal/platform/response"
	"bobshop/internal/platform/security"
	"bobshop/internal/platform/web"
)

func AuthMiddleware(parser security.Tokenizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := c.Cookie(web.AccessTokenCookieName)
		claims, err := parser.ParseToken(token)
		if err != nil {
			response.Unauthorized(c, err)
			return
		}
		c.Set(web.UserIDKey, claims["sub"])
		c.Set(web.RoleKey, claims["role"])
	}
}
