package deleteshop

import (
	"context"
)

// DeleteShopCommand represents the command to delete a shop
type DeleteShopCommand struct {
	ID string `json:"id" validate:"required,uuid"`
}

// DeleteShopCommandHandler handles the DeleteShopCommand
type DeleteShopCommandHandler interface {
	Handle(ctx context.Context, command DeleteShopCommand) error
}
