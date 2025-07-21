package infrastructure

import (
	"context"

	redis "github.com/redis/go-redis/v9"
)

const (
	maxRecentlyViewedProducts = 10
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (r *RedisCache) TrackRecentlyViewedProduct(ctx context.Context, key string, score float64, member string) error {
	if err := r.trimOldestViewedProducts(ctx, key, maxRecentlyViewedProducts); err != nil {
		return err
	}
	return r.client.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: member,
	}).Err()
}

func (r *RedisCache) GetRecentlyViewedProducts(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.ZRevRange(ctx, key, start, stop).Result()
}

func (r *RedisCache) trimOldestViewedProducts(ctx context.Context, key string, count uint32) error {
	return r.client.ZRemRangeByRank(ctx, key, 0, int64(-count-1)).Err()
}
