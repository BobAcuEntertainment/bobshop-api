package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"bobshop/internal/platform/config"
)

const (
	MongoURIWithAuth    = "mongodb://%s:%s@%s:%s/%s"
	MongoURIWithoutAuth = "mongodb://%s:%s/%s"
)

func ConnectMongo(cfg *config.DatabaseConfig) (*mongo.Client, func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var uri string
	if cfg.Username != "" && cfg.Password != "" {
		uri = fmt.Sprintf(MongoURIWithAuth, cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	} else {
		uri = fmt.Sprintf(MongoURIWithoutAuth, cfg.Host, cfg.Port, cfg.Name)
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
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
