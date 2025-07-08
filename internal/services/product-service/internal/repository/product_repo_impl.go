package repository

import (
	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/product-service/internal/db/sqlc"
)

type pgProductRepository struct {
	db      *postgresql_infra.PostgreSQLService
	queries *sqlc.Queries
}
