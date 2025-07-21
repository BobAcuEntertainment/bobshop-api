package product

import (
	"github.com/google/wire"

	"bobshop/internal/modules/product/application"
	"bobshop/internal/modules/product/delivery/http"
	"bobshop/internal/modules/product/domain"
	"bobshop/internal/modules/product/infrastructure"
)

var ProductSet = wire.NewSet(
	wire.Bind(new(domain.ProductRepository), new(*infrastructure.MongoProductRepository)),
	wire.Bind(new(domain.Cache), new(*infrastructure.RedisCache)),
	infrastructure.NewMongoProductRepository,
	infrastructure.NewRedisCache,
	application.NewProductService,
	http.NewProductHandler,
)
