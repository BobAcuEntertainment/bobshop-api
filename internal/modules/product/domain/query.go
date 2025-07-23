package domain

type ListFilter struct {
	Name       *string  `query:"name" validate:"omitempty"`
	Categories []string `query:"categories" validate:"omitempty,dive,required"`
	Brands     []string `query:"brands" validate:"omitempty,dive,required"`
	Vendor     *string  `query:"vendor" validate:"omitempty"`
	Tags       []string `query:"tags" validate:"omitempty,dive,required"`
	MinPrice   *uint32  `query:"min_price" validate:"omitempty"`
	MaxPrice   *uint32  `query:"max_price" validate:"omitempty"`
}

type CursorPagination struct {
	Cursor *string `query:"cursor" validate:"omitempty,base64"`
	Limit  *int    `query:"limit" validate:"omitempty,min=1,max=100"`
}

type Sort struct {
	SortBy *SortBy `query:"sort" validate:"omitempty,oneof=price_asc price_desc latest popular"`
}

type SortBy string

const (
	SortByPriceAsc  SortBy = "price_asc"
	SortByPriceDesc SortBy = "price_desc"
	SortByLatest    SortBy = "latest"
	SortByPopular   SortBy = "popular"
)
