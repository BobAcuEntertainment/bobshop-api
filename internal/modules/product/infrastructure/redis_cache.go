package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	redis "github.com/redis/go-redis/v9"
)

const (
	maxRecentlyViewed     = 10
	userRecentlyViewedKey = "user:%s:recently_viewed"
)

func buildUserRecentlyViewedKey(userID uuid.UUID) string {
	return fmt.Sprintf(userRecentlyViewedKey, userID)
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (r *RedisCache) TrackRecentlyViewed(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error {
	key := buildUserRecentlyViewedKey(userID)
	score := float64(time.Now().Unix())
	member := productID.String()

	if err := r.trimOldestViewed(ctx, key, maxRecentlyViewed); err != nil {
		return err
	}
	return r.client.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: member,
	}).Err()
}

func (r *RedisCache) GetRecentlyViewed(ctx context.Context, userID uuid.UUID) ([]string, error) {
	key := buildUserRecentlyViewedKey(userID)
	return r.client.ZRevRange(ctx, key, 0, int64(maxRecentlyViewed-1)).Result()
}

func (r *RedisCache) trimOldestViewed(ctx context.Context, key string, count uint32) error {
	return r.client.ZRemRangeByRank(ctx, key, 0, int64(-count-1)).Err()
}
