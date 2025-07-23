package domain

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	UpdateFields(ctx context.Context, productID uuid.UUID, fields bson.M) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
	List(ctx context.Context, filter *ListFilter, pagination *CursorPagination, sort SortBy) (products []*Product, nextCursor string, err error)
}
