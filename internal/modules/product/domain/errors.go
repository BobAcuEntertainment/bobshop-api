package domain

import "errors"

var (
	ErrProductNotFound     = errors.New("product not found")
	ErrInvalidProduct      = errors.New("invalid product")
	ErrReviewAlreadyExists = errors.New("review already exists")
)
