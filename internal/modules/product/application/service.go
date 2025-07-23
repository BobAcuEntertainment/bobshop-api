package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"

	"bobshop/internal/modules/product/domain"
)

const (
	userRecentlyViewedKey = "user:%s:recently_viewed"
)

func buildUserRecentlyViewedKey(userID uuid.UUID) string {
	return fmt.Sprintf(userRecentlyViewedKey, userID.String())
}

type ProductService struct {
	repo  domain.ProductRepository
	cache domain.Cache
}

func NewProductService(repo domain.ProductRepository, cache domain.Cache) *ProductService {
	return &ProductService{
		repo:  repo,
		cache: cache,
	}
}

func (s *ProductService) Create(ctx context.Context, product *domain.Product) error {
	return s.repo.Create(ctx, product)
}

func (s *ProductService) UpdatePartial(ctx context.Context, productID uuid.UUID, fields bson.M) error {
	return s.repo.UpdateFields(ctx, productID, fields)
}

func (s *ProductService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *ProductService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) List(
	ctx context.Context,
	filter *domain.ListFilter,
	pagination *domain.CursorPagination,
	sort domain.SortBy,
) ([]*domain.Product, string, error) {
	return s.repo.List(ctx, filter, pagination, sort)
}

func (s *ProductService) AddReview(
	ctx context.Context,
	review *domain.Review,
) error {
	product, err := s.repo.GetByID(ctx, review.ProductID)
	if err != nil {
		return err
	}
	product.AddReview(review)
	return s.repo.UpdateFields(ctx, product.ID, bson.M{"reviews": product.Reviews})
}

func (s *ProductService) TrackRecentlyViewedProduct(
	ctx context.Context,
	userID uuid.UUID,
	productID uuid.UUID,
) error {
	key := buildUserRecentlyViewedKey(userID)
	return s.cache.TrackRecentlyViewedProduct(ctx, key, float64(time.Now().Unix()), productID.String())
}

func (s *ProductService) GetRecentlyViewedProducts(
	ctx context.Context,
	userID uuid.UUID,
	limit int,
) ([]uuid.UUID, error) {
	key := buildUserRecentlyViewedKey(userID)
	res, err := s.cache.GetRecentlyViewedProducts(ctx, key, 0, int64(limit-1))
	if err != nil {
		return nil, err
	}
	ids := make([]uuid.UUID, len(res))
	for i, str := range res {
		ids[i], err = uuid.Parse(str)
		if err != nil {
			return nil, err
		}
	}
	return ids, nil
}
