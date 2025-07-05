package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/domain"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/dto"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/repository"
)

// ShopService defines the interface for shop business logic
type ShopService interface {
	// Create creates a new shop
	CreateShop(ctx context.Context, cmd *dto.CreateShopCommand) (*dto.CreateShopResponse, error)

	// GetByID retrieves a shop by its ID
	GetShopByID(ctx context.Context, id string) (*dto.ShopDetailsResponse, error)

	// GetByOwnerID retrieves shops by owner ID
	GetShopsByOwnerID(ctx context.Context, ownerID string) ([]*dto.ShopDetailsResponse, error)

	// Update updates an existing shop
	UpdateShop(ctx context.Context, id string, req *dto.UpdateShopRequest) (*dto.ShopDetailsResponse, error)

	// Delete deletes a shop
	DeleteShop(ctx context.Context, id string) error

	// List retrieves shops with pagination
	ListShops(ctx context.Context, limit, offset int) ([]*dto.ShopDetailsResponse, int64, error)

	// Search searches shops
	SearchShops(ctx context.Context, query string, limit, offset int) ([]*dto.ShopDetailsResponse, int64, error)

	// ActivateShop activates a shop
	ActivateShop(ctx context.Context, id string) error

	// BanShop bans a shop
	BanShop(ctx context.Context, id string) error
}

// shopService implements the ShopService interface
type shopService struct {
	shopRepo repository.ShopRepository
}

// NewShopService creates a new shop service
func NewShopService(shopRepo repository.ShopRepository) ShopService {
	return &shopService{
		shopRepo: shopRepo,
	}
}

// CreateShop creates a new shop
func (s *shopService) CreateShop(ctx context.Context, cmd *dto.CreateShopCommand) (*dto.CreateShopResponse, error) {
	// Check if owner already has a shop
	exists, err := s.shopRepo.ExistsByOwnerID(ctx, cmd.OwnerID)
	if err != nil {
		return nil, fmt.Errorf("failed to check owner existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("owner already has a shop")
	}

	// Check if email is already taken
	emailExists, err := s.shopRepo.ExistsByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if emailExists {
		return nil, fmt.Errorf("email already in use")
	}

	// Create shop domain object
	now := time.Now()
	shop := &domain.Shop{
		ID:              uuid.New(),
		OwnerID:         cmd.OwnerID,
		ShopName:        cmd.ShopName,
		AvatarURL:       cmd.AvatarURL,
		BannerURL:       cmd.BannerURL,
		ShopDescription: cmd.ShopDescription,
		AddressID:       cmd.AddressID,
		Phone:           cmd.Phone,
		Email:           cmd.Email,
		Rating:          0.0,
		ActiveAt:        nil, // Shop starts inactive until approved
		BannedAt:        nil,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// Save to repository
	if err := s.shopRepo.Create(ctx, shop); err != nil {
		return nil, fmt.Errorf("failed to create shop: %w", err)
	}

	// Return response
	return &dto.CreateShopResponse{
		ID:              shop.ID.String(),
		OwnerID:         shop.OwnerID.String(),
		ShopName:        shop.ShopName,
		AvatarURL:       shop.AvatarURL,
		BannerURL:       shop.BannerURL,
		ShopDescription: shop.ShopDescription,
		AddressID:       shop.AddressID.String(),
		Phone:           shop.Phone,
		Email:           shop.Email,
		Rating:          shop.Rating,
		Status:          getShopStatus(shop),
		CreatedAt:       shop.CreatedAt.Format(time.RFC3339),
	}, nil
}

// GetShopByID retrieves a shop by its ID
func (s *shopService) GetShopByID(ctx context.Context, id string) (*dto.ShopDetailsResponse, error) {
	shopID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid shop ID: %w", err)
	}

	shop, err := s.shopRepo.GetByID(ctx, shopID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shop: %w", err)
	}

	return convertToShopDetailsResponse(shop), nil
}

// GetShopsByOwnerID retrieves shops by owner ID
func (s *shopService) GetShopsByOwnerID(ctx context.Context, ownerID string) ([]*dto.ShopDetailsResponse, error) {
	ownerUUID, err := uuid.Parse(ownerID)
	if err != nil {
		return nil, fmt.Errorf("invalid owner ID: %w", err)
	}

	shops, err := s.shopRepo.GetByOwnerID(ctx, ownerUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shops by owner: %w", err)
	}

	var responses []*dto.ShopDetailsResponse
	for _, shop := range shops {
		responses = append(responses, convertToShopDetailsResponse(shop))
	}

	return responses, nil
}

// UpdateShop updates an existing shop
func (s *shopService) UpdateShop(ctx context.Context, id string, req *dto.UpdateShopRequest) (*dto.ShopDetailsResponse, error) {
	shopID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid shop ID: %w", err)
	}

	// Get existing shop
	shop, err := s.shopRepo.GetByID(ctx, shopID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shop: %w", err)
	}

	// Check if email is being changed and if new email is already taken
	if req.Email != nil && *req.Email != shop.Email {
		emailExists, err := s.shopRepo.ExistsByEmail(ctx, *req.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email existence: %w", err)
		}
		if emailExists {
			return nil, fmt.Errorf("email already in use")
		}
	}

	// Update shop fields
	shop.Update(req.ShopName, req.ShopDescription, req.AvatarURL, req.BannerURL, req.Phone, req.Email)

	// Save to repository
	if err := s.shopRepo.Update(ctx, shop); err != nil {
		return nil, fmt.Errorf("failed to update shop: %w", err)
	}

	return convertToShopDetailsResponse(shop), nil
}

// DeleteShop deletes a shop
func (s *shopService) DeleteShop(ctx context.Context, id string) error {
	shopID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid shop ID: %w", err)
	}

	if err := s.shopRepo.Delete(ctx, shopID); err != nil {
		return fmt.Errorf("failed to delete shop: %w", err)
	}

	return nil
}

// ListShops retrieves shops with pagination
func (s *shopService) ListShops(ctx context.Context, limit, offset int) ([]*dto.ShopDetailsResponse, int64, error) {
	shops, total, err := s.shopRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list shops: %w", err)
	}

	var responses []*dto.ShopDetailsResponse
	for _, shop := range shops {
		responses = append(responses, convertToShopDetailsResponse(shop))
	}

	return responses, total, nil
}

// SearchShops searches shops
func (s *shopService) SearchShops(ctx context.Context, query string, limit, offset int) ([]*dto.ShopDetailsResponse, int64, error) {
	shops, total, err := s.shopRepo.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search shops: %w", err)
	}

	var responses []*dto.ShopDetailsResponse
	for _, shop := range shops {
		responses = append(responses, convertToShopDetailsResponse(shop))
	}

	return responses, total, nil
}

// ActivateShop activates a shop
func (s *shopService) ActivateShop(ctx context.Context, id string) error {
	shopID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid shop ID: %w", err)
	}

	shop, err := s.shopRepo.GetByID(ctx, shopID)
	if err != nil {
		return fmt.Errorf("failed to get shop: %w", err)
	}

	shop.Activate()

	if err := s.shopRepo.Update(ctx, shop); err != nil {
		return fmt.Errorf("failed to activate shop: %w", err)
	}

	return nil
}

// BanShop bans a shop
func (s *shopService) BanShop(ctx context.Context, id string) error {
	shopID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid shop ID: %w", err)
	}

	shop, err := s.shopRepo.GetByID(ctx, shopID)
	if err != nil {
		return fmt.Errorf("failed to get shop: %w", err)
	}

	shop.Ban()

	if err := s.shopRepo.Update(ctx, shop); err != nil {
		return fmt.Errorf("failed to ban shop: %w", err)
	}

	return nil
}

// Helper functions

// convertToShopDetailsResponse converts domain shop to details response
func convertToShopDetailsResponse(shop *domain.Shop) *dto.ShopDetailsResponse {
	response := &dto.ShopDetailsResponse{
		ID:              shop.ID.String(),
		OwnerID:         shop.OwnerID.String(),
		ShopName:        shop.ShopName,
		AvatarURL:       shop.AvatarURL,
		BannerURL:       shop.BannerURL,
		ShopDescription: shop.ShopDescription,
		AddressID:       shop.AddressID.String(),
		Phone:           shop.Phone,
		Email:           shop.Email,
		Rating:          shop.Rating,
		Status:          getShopStatus(shop),
		CreatedAt:       shop.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       shop.UpdatedAt.Format(time.RFC3339),
	}

	if shop.ActiveAt != nil {
		activeAt := shop.ActiveAt.Format(time.RFC3339)
		response.ActiveAt = &activeAt
	}

	if shop.BannedAt != nil {
		bannedAt := shop.BannedAt.Format(time.RFC3339)
		response.BannedAt = &bannedAt
	}

	return response
}

// getShopStatus returns the current status of the shop
func getShopStatus(shop *domain.Shop) string {
	if shop.IsBanned() {
		return "banned"
	}
	if shop.IsActive() {
		return "active"
	}
	return "inactive"
}
