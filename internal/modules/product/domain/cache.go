package domain

import (
	"context"

	"github.com/google/uuid"
)

type Cache interface {
	TrackRecentlyViewed(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error
	GetRecentlyViewed(ctx context.Context, userID uuid.UUID) ([]string, error)
}
