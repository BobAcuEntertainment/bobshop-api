package domain

import "context"

type Cache interface {
	TrackRecentlyViewedProduct(ctx context.Context, key string, score float64, member string) error
	GetRecentlyViewedProducts(ctx context.Context, key string, start, stop int64) ([]string, error)
}
