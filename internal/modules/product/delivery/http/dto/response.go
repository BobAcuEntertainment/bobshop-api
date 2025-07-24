package dto

import (
	"github.com/google/uuid"

	"bobshop/internal/modules/product/domain"
)

type CreateResponse struct {
	ID uuid.UUID `json:"id"`
}

func ToCreateResponse(p *domain.Product) *CreateResponse {
	return &CreateResponse{
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

type ListResponse struct {
	Products   []*ProductResponse `json:"products"`
	NextCursor *string            `json:"next_cursor"`
}

func ToListResponse(ps []*domain.Product, nextCursor *string) *ListResponse {
	var products []*ProductResponse
	for _, p := range ps {
		products = append(products, ToProductResponse(p))
	}
	return &ListResponse{
		Products:   products,
		NextCursor: nextCursor,
	}
}
