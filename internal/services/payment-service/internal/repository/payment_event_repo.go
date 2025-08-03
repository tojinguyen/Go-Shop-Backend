package repository

import (
	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/db/sqlc"
)

type PaymentEventRepository interface {
}

type paymentEventRepository struct {
	db      *postgresql_infra.PostgreSQLService
	queries *sqlc.Queries
}

func NewPaymentEventRepository(db *postgresql_infra.PostgreSQLService) PaymentEventRepository {
	if db == nil {
		return nil
	}

	queries := sqlc.New(db.GetPool())

	return &paymentEventRepository{
		db:      db,
		queries: queries,
	}
}
