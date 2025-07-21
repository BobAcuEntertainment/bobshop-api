package dto

import (
	"github.com/google/uuid"

	"bobshop/internal/modules/product/domain"
)

type CreateProductResponse struct {
	ID uuid.UUID `json:"id"`
	// so on
}

type ProductResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	// so on
}

func FromDomain(p *domain.Product) ProductResponse {
	return ProductResponse{
		ID:   p.ID,
		Name: p.Name,
		// so on
	}
}

func FromDomainList(ps []*domain.Product) []ProductResponse {
	var res []ProductResponse
	for _, p := range ps {
		res = append(res, FromDomain(p))
	}
	return res
}
