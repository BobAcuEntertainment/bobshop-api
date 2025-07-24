package database

import (
	"context"
	"fmt"
	"log"
	"time"

	redis "github.com/redis/go-redis/v9"

	"bobshop/internal/platform/config"
)

func ConnectRedis(cfg *config.RedisConfig) (*redis.Client, func(), error) {
	opts := &redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		DB:   cfg.DB,
	}

	if cfg.Username != "" {
		opts.Username = cfg.Username
	}
	if cfg.Password != "" {
		opts.Password = cfg.Password
	}

	rdb := redis.NewClient(opts)
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
