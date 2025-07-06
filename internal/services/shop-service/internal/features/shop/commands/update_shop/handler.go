package updateshop

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/domain"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/repository"
)

// CommandHandler implements UpdateShopCommandHandler
type CommandHandler struct {
	shopRepo repository.ShopRepository
}

// NewCommandHandler creates a new UpdateShopCommandHandler
func NewCommandHandler(shopRepo repository.ShopRepository) UpdateShopCommandHandler {
	return &CommandHandler{
		shopRepo: shopRepo,
	}
}

// Handle processes the UpdateShopCommand
func (h *CommandHandler) Handle(ctx context.Context, command UpdateShopCommand) (*domain.Shop, error) {
	log.Printf("Updating shop with ID: %s", command.ID)

	existingShop, err := h.shopRepo.GetShopByID(ctx, command.ID)
	if err != nil {
		log.Printf("Error getting shop by ID: %v", err)
		return nil, fmt.Errorf("shop not found: %w", err)
	}

	shopID, err := uuid.Parse(command.ID)
	if err != nil {
		log.Printf("Error parsing shop ID: %v", err)
		return nil, fmt.Errorf("invalid shop ID: %w", err)
	}

	addressID, err := uuid.Parse(command.AddressID)
	if err != nil {
		log.Printf("Error parsing address ID: %v", err)
		return nil, fmt.Errorf("invalid address ID: %w", err)
	}

	updatedShop := &domain.Shop{
		ID:              shopID,
		OwnerID:         existingShop.OwnerID,
		ShopName:        command.ShopName,
		AvatarURL:       command.AvatarURL,
		BannerURL:       command.BannerURL,
		ShopDescription: command.ShopDescription,
		AddressID:       addressID,
		Phone:           command.Phone,
		Email:           command.Email,
		Rating:          existingShop.Rating,
		ActiveAt:        existingShop.ActiveAt,
		BannedAt:        existingShop.BannedAt,
		CreatedAt:       existingShop.CreatedAt,
	}

	err = h.shopRepo.Update(ctx, updatedShop)
	if err != nil {
		log.Printf("Error updating shop: %v", err)
		return nil, fmt.Errorf("failed to update shop: %w", err)
	}

	log.Printf("Shop updated successfully with ID: %s", command.ID)
	return updatedShop, nil
}
