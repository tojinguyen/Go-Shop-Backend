package deleteshop

import (
	"context"
	"fmt"
	"log"

	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/repository"
)

// CommandHandler implements DeleteShopCommandHandler
type CommandHandler struct {
	shopRepo repository.ShopRepository
}

// NewCommandHandler creates a new DeleteShopCommandHandler
func NewCommandHandler(shopRepo repository.ShopRepository) DeleteShopCommandHandler {
	return &CommandHandler{
		shopRepo: shopRepo,
	}
}

// Handle processes the DeleteShopCommand
func (h *CommandHandler) Handle(ctx context.Context, command DeleteShopCommand) error {
	log.Printf("Deleting shop with ID: %s", command.ID)

	// First, check if the shop exists
	_, err := h.shopRepo.GetShopByID(ctx, command.ID)
	if err != nil {
		log.Printf("Error getting shop by ID: %v", err)
		return fmt.Errorf("shop not found: %w", err)
	}

	// Delete the shop
	err = h.shopRepo.Delete(ctx, command.ID)
	if err != nil {
		log.Printf("Error deleting shop: %v", err)
		return fmt.Errorf("failed to delete shop: %w", err)
	}

	log.Printf("Shop deleted successfully with ID: %s", command.ID)
	return nil
}
