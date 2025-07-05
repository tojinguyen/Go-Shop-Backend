package createshop

import (
	"context"

	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
)

type Handler struct {
	db *postgresql_infra.PostgreSQLService
}

func (h *Handler) Handle(ctx context.Context, cmd CreateShopCommand) error {
	return nil
}
