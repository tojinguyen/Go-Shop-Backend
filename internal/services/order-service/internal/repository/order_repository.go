package repository

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/pkg/apperror"
	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/domain"
)

type OrderRepository interface {
	CreateOrder(ctx *gin.Context, order *domain.Order) (*domain.Order, error)
}

type orderRepository struct {
	db      *postgresql_infra.PostgreSQLService
	queries *sqlc.Queries
}

func NewOrderRepository(db *postgresql_infra.PostgreSQLService) OrderRepository {
	if db == nil {
		return nil
	}

	queries := sqlc.New(db.GetPool())

	return &orderRepository{
		db:      db,
		queries: queries,
	}
}

func (r *orderRepository) CreateOrder(ctx *gin.Context, order *domain.Order) (*domain.Order, error) {
	if order == nil {
		return nil, apperror.NewValidation("Order cannot be nil", fmt.Errorf("order is nil"))
	}

	//TODO: Refactor ORDER table, domain and implement this method properly

	return nil, nil
}
