package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID             uuid.UUID  `bson:"_id" json:"id"`
	CreatedAt      time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt      *time.Time `bson:"deleted_at" json:"deleted_at" default:"null"`
	IsActive       bool       `bson:"is_active" json:"is_active" default:"true"`
	Stars          [5]uint32  `bson:"stars" json:"stars" default:"[0, 0, 0, 0, 0]"`
	Sales          uint32     `bson:"sales" json:"sales" default:"0"`
	ImageGallery   []string   `bson:"image_gallery" json:"image_gallery" default:"[]"`
	ImagePrimary   string     `bson:"image_primary" json:"image_primary" default:""`
	ImageThumbnail string     `bson:"image_thumbnail" json:"image_thumbnail" default:""`
	ImageBanner    string     `bson:"image_banner" json:"image_banner" default:""`
	Brands         string     `bson:"brands" json:"brands" default:""`
	Vendor         string     `bson:"vendor" json:"vendor" default:""`
	Code           string     `bson:"code" json:"code" default:""`
	Barcode        string     `bson:"barcode" json:"barcode" default:""`
	SKU            string     `bson:"sku" json:"sku" default:""`
	Grams          uint32     `bson:"grams" json:"grams" default:"0"`
	Stock          uint32     `bson:"stock" json:"stock" default:"0"`
	PriceOld       uint32     `bson:"price_old" json:"price_old" default:"0"`
	PriceDiscount  uint32     `bson:"price_discount" json:"price_discount" default:"0"`
	Price          uint32     `bson:"price" json:"price" default:"0"`
	Content        string     `bson:"content" json:"content" default:""`
	Tags           []string   `bson:"tags" json:"tags" default:"[]"`
	Desc           string     `bson:"desc" json:"desc" default:""`
	SeoMeta        string     `bson:"seometa" json:"seometa" default:""`
	SeoTitle       string     `bson:"seotitle" json:"seotitle" default:""`
	Slug           string     `bson:"slug" json:"slug" default:""`
	NameEng        string     `bson:"nameEng" json:"nameEng" default:""`
	Name           string     `bson:"name" json:"name" default:""`
	Reviews        []*Review  `bson:"reviews" json:"reviews" default:"[]"`
	Categories     []string   `bson:"categories" json:"categories" default:"[]"`
}

func (p *Product) AddReview(review *Review) {
	p.Reviews = append(p.Reviews, review)
	p.Stars[review.Rating-1]++
}
