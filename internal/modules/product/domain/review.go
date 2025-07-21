package domain

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID        uuid.UUID `bson:"_id" json:"id"`
	ProductID uuid.UUID `bson:"product_id" json:"product_id"`
	UserID    uuid.UUID `bson:"user_id" json:"user_id"`
	Rating    uint8     `bson:"rating" json:"rating" validate:"min=1,max=5"`
	Comment   string    `bson:"comment" json:"comment"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func NewReview(productID, userID uuid.UUID, rating uint8, comment string) *Review {
	return &Review{
		ID:        uuid.New(),
		ProductID: productID,
		UserID:    userID,
		Rating:    rating,
		Comment:   comment,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
