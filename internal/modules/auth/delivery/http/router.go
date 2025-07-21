package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *AuthHandler) {
	authRoutes := rg.Group("/auth")
	{
		authRoutes.POST("/signup", handler.SignUp)
		authRoutes.POST("/signin", handler.SignIn)
		authRoutes.POST("/signout", handler.SignOut)
	}
}
