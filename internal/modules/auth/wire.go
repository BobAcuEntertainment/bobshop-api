package auth

import (
	"github.com/google/wire"

	"bobshop/internal/modules/auth/application"
	"bobshop/internal/modules/auth/delivery/http"
	"bobshop/internal/modules/auth/domain"
	"bobshop/internal/modules/auth/infrastructure"
)

var AuthSet = wire.NewSet(
	wire.Bind(new(domain.AuthRepository), new(*infrastructure.MongoAuthRepository)),
	infrastructure.NewMongoAuthRepository,
	application.NewAuthService,
	http.NewAuthHandler,
)
