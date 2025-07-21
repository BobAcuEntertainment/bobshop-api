package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"bobshop/internal/platform/config"
	"bobshop/internal/platform/response"

	authHttp "bobshop/internal/modules/auth/delivery/http"
	productHttp "bobshop/internal/modules/product/delivery/http"
)

type AppServer struct {
	Engine *gin.Engine
}

func provideGinEngine(cfg *config.ServerConfig) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	if len(cfg.TrustedProxies) == 0 {
		engine.SetTrustedProxies(nil)
	} else {
		engine.SetTrustedProxies(cfg.TrustedProxies)
	}
	switch cfg.Mode {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "test":
		gin.SetMode(gin.TestMode)
	case "release":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
	return engine
}

func initializeServer(
	engine *gin.Engine,
	authMiddleware gin.HandlerFunc,
	authHandler *authHttp.AuthHandler,
	productHandler *productHttp.ProductHandler,
) *AppServer {
	// Register global middleware here if any

	// Register routes
	apiV1 := engine.Group("/api/v1")

	// Health check endpoint
	apiV1.GET("/health", func(c *gin.Context) {
		response.Success(c, http.StatusOK, "OK", nil)
	})

	// auth routes
	authHttp.RegisterRoutes(apiV1, authHandler)

	// product routes
	productHttp.RegisterRoutes(apiV1, authMiddleware, productHandler)

	return &AppServer{
		Engine: engine,
	}
}
