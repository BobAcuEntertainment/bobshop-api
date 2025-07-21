package database

import (
	"context"
	"log"
	"time"

	redis "github.com/redis/go-redis/v9"

	"bobshop/internal/platform/config"
)

func ConnectRedis(cfg *config.RedisConfig) (*redis.Client, func(), error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		if err := rdb.Close(); err != nil {
			log.Printf("Failed to close Redis connection: %v", err)
		}
	}

	log.Println("Successfully connected to Redis.")
	return rdb, cleanup, nil
}
