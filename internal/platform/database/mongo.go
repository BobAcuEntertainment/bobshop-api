package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"bobshop/internal/platform/config"
)

func ConnectMongo(cfg *config.DatabaseConfig) (*mongo.Client, func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Failed to disconnect from MongoDB: %v", err)
		}
	}

	log.Println("Successfully connected to MongoDB.")
	return client, cleanup, nil
}

func ProvideMongoDatabase(client *mongo.Client, cfg *config.DatabaseConfig) *mongo.Database {
	return client.Database(cfg.Name)
}
