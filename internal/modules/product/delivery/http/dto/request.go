package dto

import (
	"github.com/google/uuid"

	"bobshop/internal/modules/product/domain"
)

type CreateProductRequest struct {
	Name string `json:"name" validate:"required"`
	// so on
}

func (r *CreateProductRequest) ToDomain() *domain.Product {
	return &domain.Product{
		Name: r.Name,
		// so on
	}
}

type UpdateProductRequest struct {
	Name string `json:"name" validate:"omitempty"`
	// ...
}

func (r *UpdateProductRequest) ToDomain(id uuid.UUID) *domain.Product {
	return &domain.Product{
		ID:   id,
		Name: r.Name,
		// so on
	}
}

type AddReviewRequest struct {
	Rating  uint8  `json:"rating" validate:"required,min=1,max=5"`
	Comment string `json:"comment" validate:"omitempty"`
}

func (r *AddReviewRequest) ToDomain(productID, userID uuid.UUID) *domain.Review {
	return domain.NewReview(productID, userID, r.Rating, r.Comment)
}
