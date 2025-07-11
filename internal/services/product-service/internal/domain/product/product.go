package product

import (
	"errors"
	"time"

	time_utils "github.com/toji-dev/go-shop/internal/pkg/time"

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

func Reconstitute(
	id, shopID, categoryID uuid.UUID,
	name, description, thumbnailURL string,
	price Price,
	quantity int,
	status ProductStatus,
	createdAt, updatedAt time.Time,
) (*Product, error) {
	return &Product{
		id:           id,
		shopID:       shopID,
		name:         name,
		description:  description,
		categoryID:   categoryID,
		price:        price,
		quantity:     quantity,
		thumbnailURL: thumbnailURL,
		status:       status,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}, nil
}

func (p *Product) ChangeName(newName string) error {
	if newName == "" {
		return errors.New("product name cannot be empty")
	}
	p.name = newName
	p.setUpdatedAt()
	return nil
}

func (p *Product) UpdateDescription(newDescription string) {
	p.description = newDescription
	p.setUpdatedAt()
}

func (p *Product) UpdateThumbnail(newURL string) error {
	p.thumbnailURL = newURL
	p.setUpdatedAt()
	return nil
}

func (p *Product) ChangeCategory(newCategoryID uuid.UUID) {
	p.categoryID = newCategoryID
	p.setUpdatedAt()
}

func (p *Product) ChangePrice(newPrice Price) error {
	if newPrice.GetAmount() < 0 {
		return errors.New("price cannot be negative")
	}
	p.price = newPrice
	p.setUpdatedAt()
	return nil
}

func (p *Product) UpdateQuantity(newQuantity int) error {
	if newQuantity < 0 {
		return errors.New("quantity cannot be negative")
	}
	p.quantity = newQuantity
	if newQuantity == 0 {
		p.status = ProductStatusOutOfStock
	} else {
		if p.status == ProductStatusOutOfStock {
			p.status = ProductStatusActive
		}
	}
	p.setUpdatedAt()
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

func (p *Product) setUpdatedAt() {
	p.updatedAt = time_utils.GetUtcTime()
}
