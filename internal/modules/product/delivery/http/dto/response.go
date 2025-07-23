package dto

import (
	"github.com/google/uuid"

	"bobshop/internal/modules/product/domain"
)

type CreateProductResponse struct {
	ID uuid.UUID `json:"id"`
}

func ToCreateProductResponse(p *domain.Product) *CreateProductResponse {
	return &CreateProductResponse{
		ID: p.ID,
	}
}

type ProductResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func ToProductResponse(p *domain.Product) *ProductResponse {
	return &ProductResponse{
		ID:   p.ID,
		Name: p.Name,
	}
}

type ListProductsResponse struct {
	Products   []*ProductResponse `json:"products"`
	NextCursor *string            `json:"next_cursor"`
}

func ToListProductsResponse(ps []*domain.Product, nextCursor *string) *ListProductsResponse {
	var products []*ProductResponse
	for _, p := range ps {
		products = append(products, ToProductResponse(p))
	}
	return &ListProductsResponse{
		Products:   products,
		NextCursor: nextCursor,
	}
}
