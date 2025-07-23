package domain

import (
	"context"

	"github.com/google/uuid"
)

type Cache interface {
	TrackRecentlyViewedProduct(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error
	GetRecentlyViewedProducts(ctx context.Context, userID uuid.UUID, limit int) ([]string, error)
}
