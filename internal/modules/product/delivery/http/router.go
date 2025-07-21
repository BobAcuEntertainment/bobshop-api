package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(group *gin.RouterGroup, authMiddleware gin.HandlerFunc, h *ProductHandler) {
	products := group.Group("/products")
	{
		admin := products.Group("", authMiddleware)
		admin.POST("", h.Create)
		admin.PUT("/:id", h.Update)
		admin.DELETE("/:id", h.Delete)

		products.GET("/:id", h.GetByID)
		products.GET("", h.List)
		products.POST("/:id/reviews", h.AddReview)
		products.POST("/:id/view", h.TrackRecentlyViewed)
		products.GET("/recently-viewed", h.GetRecentlyViewed)
	}
}
