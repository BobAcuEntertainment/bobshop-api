//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"bobshop/internal/modules/auth"
	"bobshop/internal/modules/product"
	"bobshop/internal/platform/config"
	"bobshop/internal/platform/database"
	"bobshop/internal/platform/infrastructure"
	"bobshop/internal/platform/middleware"
	"bobshop/internal/platform/security"
	"bobshop/internal/platform/web"
)

func buildApp(cfg *config.Config) (*AppServer, func(), error) {
	panic(wire.Build(
		// Config
		wire.FieldsOf(new(*config.Config), "Server", "Database", "Redis", "Cookie", "JWT"),
		wire.Bind(new(security.Tokenizer), new(*infrastructure.JwtTokenizer)),

		// Infrastructure
		database.ConnectMongo,
		database.ProvideMongoDatabase,
		database.ConnectRedis,

		// JWT
		infrastructure.NewJwtTokenizer,

		// Cookie
		web.NewCookieManager,

		// Middlewares
		middleware.AuthMiddleware,

		// Modules
		auth.AuthSet,
		product.ProductSet,

		// Presentation
		provideGinEngine,
		initializeServer,
	))
}
