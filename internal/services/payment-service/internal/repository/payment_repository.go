package repository

import (
	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
)

type PaymentRepository interface {
}

type paymentRepository struct {
	db *postgresql_infra.PostgreSQLService
	// queries *sqlc.Queries
}

func NewPaymentRepository(db *postgresql_infra.PostgreSQLService) PaymentRepository {
	if db == nil {
		return nil
	}

	// queries := sqlc.New(db.GetPool())

	return &paymentRepository{
		db: db,
		// queries: queries,
	}
}
