package repository

import (
	"context"

	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/domain"
)

type pgProductRepository struct {
	db      *postgresql_infra.PostgreSQLService
	queries *sqlc.Queries
}

func NewProductRepository(db *postgresql_infra.PostgreSQLService) ProductRepository {
	if db == nil {
		return nil
	}

	queries := sqlc.New(db.GetPool())
	return &pgProductRepository{
		db:      db,
		queries: queries,
	}
}

func (r *pgProductRepository) Save(ctx context.Context, product *domain.Product) error {

	return nil
}

func (r *pgProductRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	return nil, nil
}

func (r *pgProductRepository) GetByShopID(ctx context.Context, shopID string) ([]*domain.Product, error) {
	return nil, nil
}
