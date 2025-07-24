package dto

type CreateRequest struct {
	Name  string `json:"name" validate:"required"`
	Price uint32 `json:"price" validate:"required"`
}

type UpdateRequest struct {
	Name  *string `json:"name" validate:"omitempty"`
	Price *uint32 `json:"price" validate:"omitempty"`
}

type AddReviewRequest struct {
	Rating  uint8  `json:"rating" validate:"required,min=1,max=5"`
	Comment string `json:"comment" validate:"omitempty"`
}

type ListFilterRequest struct {
	Name       *string  `query:"name" validate:"omitempty"`
	Categories []string `query:"categories" validate:"omitempty,dive,required"`
	Brands     []string `query:"brands" validate:"omitempty,dive,required"`
	Vendor     *string  `query:"vendor" validate:"omitempty"`
	Tags       []string `query:"tags" validate:"omitempty,dive,required"`
	MinPrice   *uint32  `query:"min_price" validate:"omitempty"`
	MaxPrice   *uint32  `query:"max_price" validate:"omitempty"`
}

type CursorPaginationRequest struct {
	Cursor *string `query:"cursor" validate:"omitempty,base64"`
	Limit  *int    `query:"limit" validate:"omitempty,min=1,max=100"`
}

type SortRequest struct {
	SortBy *string `query:"sort" validate:"omitempty,oneof=price_asc price_desc latest popular"`
}
