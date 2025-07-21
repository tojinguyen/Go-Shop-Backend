package repository

import (
	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
)

type OrderRepository interface {
}

type orderRepository struct {
	db *postgresql_infra.PostgreSQLService
}

func NewOrderRepository(db *postgresql_infra.PostgreSQLService) OrderRepository {
	return &orderRepository{
		db: db,
	}
}
