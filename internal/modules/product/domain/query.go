package domain

type ListFilter struct {
	Name       *string  `validate:"omitempty"`
	Categories []string `validate:"omitempty,dive,required"`
	Brands     []string `validate:"omitempty,dive,required"`
	Vendor     *string  `validate:"omitempty"`
	Tags       []string `validate:"omitempty,dive,required"`
	MinPrice   *uint32  `validate:"omitempty"`
	MaxPrice   *uint32  `validate:"omitempty"`
}

type CursorPagination struct {
	Cursor *string `validate:"omitempty,base64"`
	Limit  *int    `validate:"omitempty,min=1,max=100"`
}

type Sort struct {
	SortBy *SortBy `validate:"omitempty,oneof=price_asc price_desc latest popular"`
}

type SortBy string

const (
	SortByPriceAsc  SortBy = "price_asc"
	SortByPriceDesc SortBy = "price_desc"
	SortByLatest    SortBy = "latest"
	SortByPopular   SortBy = "popular"
)
