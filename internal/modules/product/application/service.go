package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"

	"bobshop/internal/modules/product/delivery/http/dto"
	"bobshop/internal/modules/product/domain"
)

func fromCreateRequest(req dto.CreateRequest) *domain.Product {
	return domain.NewProductBuilder(req.Name, req.Price).Build()
}

func fromUpdateRequest(req dto.UpdateRequest) bson.M {
	updateFields := bson.M{
		"updated_at": time.Now(),
	}
	if req.Name != nil {
		updateFields["name"] = *req.Name
	}
	if req.Price != nil {
		updateFields["price"] = *req.Price
	}
	return updateFields
}

func fromAddReviewRequest(req dto.AddReviewRequest, productID, userID uuid.UUID) *domain.Review {
	return domain.NewReview(productID, userID, req.Rating, req.Comment)
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

func (s *ProductService) Create(ctx context.Context, req dto.CreateRequest) (*domain.Product, error) {
	product := fromCreateRequest(req)
	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) Update(ctx context.Context, productID uuid.UUID, req dto.UpdateRequest) error {
	updateFields := fromUpdateRequest(req)
	return s.repo.Update(ctx, productID, updateFields)
}

func (s *ProductService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *ProductService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func fromListFilterRequest(req dto.ListFilterRequest) *domain.ListFilter {
	return &domain.ListFilter{
		Name:       req.Name,
		Categories: req.Categories,
		Brands:     req.Brands,
		Vendor:     req.Vendor,
		Tags:       req.Tags,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
	}
}

func fromCursorPaginationRequest(req dto.CursorPaginationRequest) *domain.CursorPagination {
	return &domain.CursorPagination{
		Cursor: req.Cursor,
		Limit:  req.Limit,
	}
}

func fromSortRequest(req dto.SortRequest) *domain.Sort {
	sortBy := domain.SortBy(*req.SortBy)
	return &domain.Sort{
		SortBy: &sortBy,
	}
}

func (s *ProductService) List(
	ctx context.Context,
	filterRequest dto.ListFilterRequest,
	paginationRequest dto.CursorPaginationRequest,
	sortRequest dto.SortRequest,
) ([]*domain.Product, *string, error) {
	filter := fromListFilterRequest(filterRequest)
	pagination := fromCursorPaginationRequest(paginationRequest)
	sort := fromSortRequest(sortRequest)
	return s.repo.List(ctx, filter, pagination, sort)
}

func (s *ProductService) AddReview(
	ctx context.Context,
	req dto.AddReviewRequest,
	userID uuid.UUID,
	productID uuid.UUID,
) error {
	review := fromAddReviewRequest(req, productID, userID)
	return s.repo.AddReview(ctx, review)
}

func (s *ProductService) TrackRecentlyViewed(
	ctx context.Context,
	userID uuid.UUID,
	productID uuid.UUID,
) error {
	return s.cache.TrackRecentlyViewed(ctx, userID, productID)
}

func (s *ProductService) GetRecentlyViewed(
	ctx context.Context,
	userID uuid.UUID,
) ([]uuid.UUID, error) {
	res, err := s.cache.GetRecentlyViewed(ctx, userID)
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
