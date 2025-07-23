package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID             uuid.UUID  `bson:"_id" json:"id" validate:"required,uuid4"`
	CreatedAt      time.Time  `bson:"created_at" json:"created_at" validate:"required"`
	UpdatedAt      time.Time  `bson:"updated_at" json:"updated_at" validate:"required"`
	DeletedAt      *time.Time `bson:"deleted_at" json:"deleted_at"`
	IsActive       bool       `bson:"is_active" json:"is_active"`
	Stars          [5]uint32  `bson:"stars" json:"stars"`
	Sales          uint32     `bson:"sales" json:"sales"`
	ImageGallery   []string   `bson:"image_gallery" json:"image_gallery" validate:"dive,url"`
	ImagePrimary   string     `bson:"image_primary" json:"image_primary" validate:"omitempty,url"`
	ImageThumbnail string     `bson:"image_thumbnail" json:"image_thumbnail" validate:"omitempty,url"`
	ImageBanner    string     `bson:"image_banner" json:"image_banner" validate:"omitempty,url"`
	Brands         string     `bson:"brands" json:"brands" validate:"omitempty"`
	Vendor         string     `bson:"vendor" json:"vendor" validate:"omitempty"`
	Code           string     `bson:"code" json:"code" validate:"omitempty"`
	Barcode        string     `bson:"barcode" json:"barcode" validate:"omitempty"`
	SKU            string     `bson:"sku" json:"sku" validate:"omitempty"`
	Grams          uint32     `bson:"grams" json:"grams" validate:"omitempty"`
	Stock          uint32     `bson:"stock" json:"stock" validate:"omitempty"`
	PriceOld       uint32     `bson:"price_old" json:"price_old" validate:"omitempty"`
	PriceDiscount  uint32     `bson:"price_discount" json:"price_discount" validate:"omitempty"`
	Price          uint32     `bson:"price" json:"price" validate:"required"`
	Content        string     `bson:"content" json:"content" validate:"omitempty"`
	Tags           []string   `bson:"tags" json:"tags" validate:"dive,required"`
	Desc           string     `bson:"desc" json:"desc" validate:"omitempty"`
	SeoMeta        string     `bson:"seometa" json:"seometa" validate:"omitempty"`
	SeoTitle       string     `bson:"seotitle" json:"seotitle" validate:"omitempty"`
	Slug           string     `bson:"slug" json:"slug" validate:"omitempty"`
	NameEng        string     `bson:"nameEng" json:"nameEng" validate:"omitempty"`
	Name           string     `bson:"name" json:"name" validate:"required"`
	Reviews        []*Review  `bson:"reviews" json:"reviews" validate:"dive"`
	Categories     []string   `bson:"categories" json:"categories" validate:"dive,required"`
}

type ProductBuilder struct {
	product *Product
}

func NewProductBuilder(name string, price uint32) *ProductBuilder {
	return &ProductBuilder{
		product: &Product{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsActive:  true,
			Name:      name,
			Price:     price,
		},
	}
}

func (b *ProductBuilder) WithDesc(desc string) *ProductBuilder {
	b.product.Desc = desc
	return b
}

func (b *ProductBuilder) Build() *Product {
	return b.product
}
