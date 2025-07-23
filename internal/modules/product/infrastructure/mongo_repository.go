package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"bobshop/internal/modules/product/domain"
)

type MongoProductRepository struct {
	collection *mongo.Collection
}

func NewMongoProductRepository(db *mongo.Database) *MongoProductRepository {
	collection := db.Collection("products")

	// Create indexes in migrations, also consider using compound indexes
	// indexModels := []mongo.IndexModel{
	// 	{Keys: bson.D{{"slug", 1}}, Options: options.Index().SetUnique(true)},
	// 	{Keys: bson.D{{"categories", 1}}},
	// 	{Keys: bson.D{{"brandlist", 1}}},
	// 	{Keys: bson.D{{"vendor", 1}}},
	// 	{Keys: bson.D{{"tags", 1}}},
	// 	{Keys: bson.D{{"price", 1}}},
	// 	{Keys: bson.D{{"created_at", -1}}},
	// 	{Keys: bson.D{{"sale", -1}}},
	// }
	// _, err := collection.Indexes().CreateMany(context.Background(), indexModels)
	// if err != nil {
	// 	return nil, err
	// }

	return &MongoProductRepository{collection: collection}
}

func (r *MongoProductRepository) Create(ctx context.Context, product *domain.Product) error {
	_, err := r.collection.InsertOne(ctx, product)
	return err
}

func (r *MongoProductRepository) UpdateFields(ctx context.Context, productID uuid.UUID, fields bson.M) error {
	filter := bson.M{"_id": productID}
	update := bson.M{"$set": fields}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

func (r *MongoProductRepository) AddReview(ctx context.Context, review *domain.Review) error {
	filter := bson.M{"user_id": review.UserID, "product_id": review.ProductID}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
			return err
	}
	if count > 0 {
			return domain.ErrReviewAlreadyExists
	}
	
	update := bson.M{
			"$push": bson.M{"reviews": review},
			"$inc": bson.M{fmt.Sprintf("stars.%d", review.Rating-1): 1},
	}

	return r.UpdateFields(ctx, review.ProductID, update)
}

func (r *MongoProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	filter := bson.M{"_id": id, "deleted_at": nil}
	update := bson.M{"$set": bson.M{"deleted_at": time.Now(), "is_active": false}}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

func (r *MongoProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var product domain.Product
	filter := bson.M{"_id": id, "deleted_at": nil, "is_active": true}
	err := r.collection.FindOne(ctx, filter).Decode(&product)
	if err == mongo.ErrNoDocuments {
		return nil, domain.ErrProductNotFound
	}
	return &product, err
}

func (r *MongoProductRepository) List(
	ctx context.Context,
	filter *domain.ListFilter,
	pagination *domain.CursorPagination,
	sort *domain.Sort,
) ([]*domain.Product, *string, error) {
	query := bson.M{"deleted_at": nil, "is_active": true}

	if filter != nil {
		if filter.Name != nil {
			query["name"] = bson.M{"$regex": *filter.Name, "$options": "i"}
		}
		if len(filter.Categories) > 0 {
			query["categories"] = bson.M{"$in": filter.Categories}
		}
		if len(filter.Brands) > 0 {
			query["brandlist"] = bson.M{"$in": filter.Brands}
		}
		if filter.Vendor != nil {
			query["vendor"] = *filter.Vendor
		}
		if len(filter.Tags) > 0 {
			query["tags"] = bson.M{"$in": filter.Tags}
		}
		if filter.MinPrice != nil {
			query["price"] = bson.M{"$gte": *filter.MinPrice}
		}
		if filter.MaxPrice != nil {
			if existing, ok := query["price"].(bson.M); ok {
				existing["$lte"] = *filter.MaxPrice
			} else {
				query["price"] = bson.M{"$lte": *filter.MaxPrice}
			}
		}
	}

	opts := options.Find().SetLimit(int64(*pagination.Limit + 1)) // +1 for next cursor

	var sortBson bson.D
	switch *sort.SortBy {
	case domain.SortByPriceAsc:
		sortBson = bson.D{{Key: "price", Value: 1}, {Key: "_id", Value: 1}}
	case domain.SortByPriceDesc:
		sortBson = bson.D{{Key: "price", Value: -1}, {Key: "_id", Value: 1}}
	case domain.SortByLatest:
		sortBson = bson.D{{Key: "created_at", Value: -1}, {Key: "_id", Value: 1}}
	case domain.SortByPopular:
		sortBson = bson.D{{Key: "sale", Value: -1}, {Key: "_id", Value: 1}} // Assuming popular by sales
	default:
		sortBson = bson.D{{Key: "_id", Value: 1}}
	}
	opts.SetSort(sortBson)

	if pagination.Cursor != nil {
		cursorID, err := uuid.Parse(*pagination.Cursor)
		if err != nil {
			return nil, nil, err
		}
		// For cursor, assume cursor is last _id, and sort includes _id
		query["_id"] = bson.M{"$gte": cursorID}
	}

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	var products []*domain.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, nil, err
	}

	var nextCursor string
	if len(products) > *pagination.Limit {
		nextCursor = products[*pagination.Limit].ID.String()
		products = products[:*pagination.Limit]
	}

	return products, &nextCursor, nil
}
