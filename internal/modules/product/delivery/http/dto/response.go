package dto

import (
	"github.com/google/uuid"

	"bobshop/internal/modules/product/domain"
)

type CreateProductResponse struct {
	ID uuid.UUID `json:"id"`
}

func ToResponse(p *domain.Product) *CreateProductResponse {
	return &CreateProductResponse{
		ID: p.ID,
	}
}

type ProductResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func FromDomain(p *domain.Product) ProductResponse {
	return ProductResponse{
		ID:   p.ID,
		Name: p.Name,
	}
}

func FromDomainList(ps []*domain.Product) []ProductResponse {
	var res []ProductResponse
	for _, p := range ps {
		res = append(res, FromDomain(p))
	}
	return res
}
