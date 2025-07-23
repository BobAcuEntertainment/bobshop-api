package dto

import (
	"github.com/google/uuid"

	"bobshop/internal/modules/product/domain"
)

type CreateProductRequest struct {
	Name  string `json:"name" validate:"required"`
	Price uint32 `json:"price" validate:"required"`
}

func ToDomain(r *CreateProductRequest) *domain.Product {
	return domain.NewProductBuilder(r.Name, r.Price).Build()
}

type UpdateProductRequest struct {
	Name  *string `json:"name"`
	Price *uint32 `json:"price"`
}

type AddReviewRequest struct {
	Rating  uint8  `json:"rating" validate:"required,min=1,max=5"`
	Comment string `json:"comment"`
}

func (r *AddReviewRequest) ToDomain(productID, userID uuid.UUID) *domain.Review {
	return domain.NewReview(productID, userID, r.Rating, r.Comment)
}
