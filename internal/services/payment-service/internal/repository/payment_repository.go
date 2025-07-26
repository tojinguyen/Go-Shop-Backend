package repository

import (
	"context"

	"github.com/toji-dev/go-shop/internal/pkg/converter"
	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/domain"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, params sqlc.CreatePaymentParams) (*domain.Payment, error)
	UpdatePaymentStatus(ctx context.Context, params sqlc.UpdatePaymentStatusParams) (*domain.Payment, error)
	GetPaymentByOrderID(ctx context.Context, orderID string) (*domain.Payment, error)
}

type paymentRepository struct {
	db      *postgresql_infra.PostgreSQLService
	queries *sqlc.Queries
}

func NewPaymentRepository(db *postgresql_infra.PostgreSQLService) PaymentRepository {
	if db == nil {
		return nil
	}

	queries := sqlc.New(db.GetPool())

	return &paymentRepository{
		db:      db,
		queries: queries,
	}
}

func (r *paymentRepository) CreatePayment(ctx context.Context, params sqlc.CreatePaymentParams) (*domain.Payment, error) {
	result, err := r.queries.CreatePayment(ctx, params)
	if err != nil {
		return nil, err
	}
	return toDomain(&result), nil
}

func (r *paymentRepository) UpdatePaymentStatus(ctx context.Context, params sqlc.UpdatePaymentStatusParams) (*domain.Payment, error) {
	result, err := r.queries.UpdatePaymentStatus(ctx, params)
	if err != nil {
		return nil, err
	}
	return toDomain(&result), nil
}

func (r *paymentRepository) GetPaymentByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	result, err := r.queries.GetPaymentByOrderID(ctx, converter.StringToPgUUID(orderID))
	if err != nil {
		return nil, err
	}
	return toDomain(&result), nil
}

// toDomain converts sqlc.Payment to domain.Payment
func toDomain(p *sqlc.Payment) *domain.Payment {
	return &domain.Payment{
		ID:                    converter.PgUUIDToString(p.ID),
		OrderID:               converter.PgUUIDToString(p.OrderID),
		UserID:                converter.PgUUIDToString(p.UserID),
		Amount:                converter.PgNumericToFloat64(p.Amount),
		Currency:              p.Currency,
		Method:                domain.PaymentMethod(p.PaymentMethod),
		Provider:              *converter.PgTextToStringPtr(p.PaymentProvider),
		ProviderTransactionID: converter.PgTextToStringPtr(p.ProviderTransactionID),
		Status:                domain.PaymentStatus(p.PaymentStatus),
		CreatedAt:             p.CreatedAt.Time,
		UpdatedAt:             p.UpdatedAt.Time,
	}
}
