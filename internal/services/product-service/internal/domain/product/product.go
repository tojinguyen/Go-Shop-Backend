package domain

import (
	"errors"
	"time"

	"github.com/toji-dev/go-shop/internal/pkg/converter"
	time_utils "github.com/toji-dev/go-shop/internal/pkg/time"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/db/sqlc"

	uuid "github.com/google/uuid"
)

type ProductStatus string

const (
	ProductStatusActive       ProductStatus = "ACTIVE"
	ProductStatusInactive     ProductStatus = "INACTIVE"
	ProductStatusOutOfStock   ProductStatus = "OUT_OF_STOCK"
	ProductStatusDiscontinued ProductStatus = "DISCONTINUED"
	ProductStatusBanned       ProductStatus = "BANNED"
)

type Product struct {
	id              uuid.UUID
	shopID          uuid.UUID
	name            string
	thumbnailURL    string
	description     string
	categoryID      uuid.UUID
	price           Price
	quantity        int
	reserveQuantity int
	status          ProductStatus
	soldCount       int
	ratingAvg       float64
	totalReviews    int
	createdAt       time.Time
	updatedAt       time.Time
	deletedAt       *time.Time // Nullable field for soft deletion
}

func NewProduct(shopID, name, thumbnailURL, description string, categoryID uuid.UUID, price Price, quantity int) (*Product, error) {
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}
	return &Product{
		id:              uuid.New(),
		shopID:          uuid.MustParse(shopID),
		name:            name,
		thumbnailURL:    thumbnailURL,
		description:     description,
		categoryID:      categoryID,
		price:           price,
		quantity:        quantity,
		reserveQuantity: quantity,
		status:          ProductStatusActive,
		soldCount:       0,
		ratingAvg:       0.0,
		totalReviews:    0,
		createdAt:       time_utils.GetUtcTime(),
		updatedAt:       time_utils.GetUtcTime(),
		deletedAt:       nil,
	}, nil
}

func ConvertFromSqlcProduct(sqlcProduct *sqlc.Product) (*Product, error) {
	price, err := NewPrice(
		*converter.PgNumericToFloat64Ptr(sqlcProduct.Price),
		"USD",
	)
	if err != nil {
		return nil, err
	}

	return &Product{
		id:              converter.PgUUIDToUUID(sqlcProduct.ID),
		shopID:          converter.PgUUIDToUUID(sqlcProduct.ShopID),
		name:            sqlcProduct.ProductName,
		thumbnailURL:    sqlcProduct.ThumbnailUrl.String,
		description:     sqlcProduct.ProductDescription.String,
		categoryID:      converter.PgUUIDToUUID(sqlcProduct.CategoryID),
		price:           price,
		quantity:        int(sqlcProduct.Quantity),
		reserveQuantity: int(sqlcProduct.ReserveQuantity),
		status:          ProductStatus(sqlcProduct.ProductStatus),
		soldCount:       int(sqlcProduct.SoldCount),
		ratingAvg:       *converter.PgNumericToFloat64Ptr(sqlcProduct.RatingAvg),
		totalReviews:    int(sqlcProduct.TotalReviews),
		createdAt:       sqlcProduct.CreatedAt.Time,
		updatedAt:       sqlcProduct.UpdatedAt.Time,
		deletedAt:       nil,
	}, nil
}

func (p *Product) ChangePrice(newPrice Price) error {
	if newPrice.GetAmount() < 0 {
		return errors.New("price cannot be negative")
	}

	p.price = newPrice
	p.updatedAt = time_utils.GetUtcTime()
	return nil
}

func (p *Product) UpdateQuantity(newQuantity int) error {
	if newQuantity < 0 {
		return errors.New("quantity cannot be negative")
	}

	p.quantity = newQuantity
	p.updatedAt = time_utils.GetUtcTime()
	return nil
}

func (p *Product) Deactivate() {
	p.status = ProductStatusInactive
	p.updatedAt = time_utils.GetUtcTime()
}

func (p *Product) ID() uuid.UUID         { return p.id }
func (p *Product) ShopID() uuid.UUID     { return p.shopID }
func (p *Product) Name() string          { return p.name }
func (p *Product) Description() *string  { return &p.description }
func (p *Product) ThumbnailURL() *string { return &p.thumbnailURL }
func (p *Product) CategoryID() uuid.UUID { return p.categoryID }
func (p *Product) Price() Price          { return p.price }
func (p *Product) Quantity() int         { return p.quantity }
func (p *Product) Status() ProductStatus { return p.status }
func (p *Product) CreatedAt() time.Time  { return p.createdAt }
func (p *Product) UpdatedAt() time.Time  { return p.updatedAt }
