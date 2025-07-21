package infrastructure

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"bobshop/internal/modules/auth/domain"
)

type MongoAuthRepository struct {
	collection *mongo.Collection
}

func NewMongoAuthRepository(db *mongo.Database) *MongoAuthRepository {
	return &MongoAuthRepository{collection: db.Collection("users")}
}

func (r *MongoAuthRepository) Create(ctx context.Context, user *domain.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *MongoAuthRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}
