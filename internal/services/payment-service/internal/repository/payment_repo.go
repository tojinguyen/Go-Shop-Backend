package repository

import (
	"context"

	"github.com/toji-dev/go-shop/internal/pkg/converter"
	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/constant"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/domain"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, params sqlc.CreatePaymentParams) (*domain.Payment, error)
	UpdatePaymentStatus(ctx context.Context, params sqlc.UpdatePaymentStatusParams) (*domain.Payment, error)
	GetPaymentByOrderID(ctx context.Context, orderID string) (*domain.Payment, error)
	CreatePaymentRefund(ctx context.Context, params sqlc.CreateRefundPaymentParams) (*domain.PaymentRefund, error)
	GetRefundByPaymentID(ctx context.Context, paymentID string) (*domain.PaymentRefund, error)
	UpdateRefundPaymentStatus(ctx context.Context, params sqlc.UpdateRefundPaymentStatusParams) (*domain.PaymentRefund, error)
	GetBatchRefundPaymentsByStatus(ctx context.Context, status sqlc.RefundStatus) ([]domain.PaymentRefund, error)
	GetBatchPendingPayments(ctx context.Context) ([]domain.Payment, error)
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

// toDomain converts sqlc.Payment to domain.Payment
func toDomain(p *sqlc.Payment) *domain.Payment {
	return &domain.Payment{
		ID:                    converter.PgUUIDToString(p.ID),
		OrderID:               converter.PgUUIDToString(p.OrderID),
		UserID:                converter.PgUUIDToString(p.UserID),
		Amount:                converter.PgNumericToFloat64(p.Amount),
		Currency:              p.Currency,
		Method:                constant.PaymentMethod(p.PaymentMethod),
		Provider:              *converter.PgTextToStringPtr(p.PaymentProvider),
		ProviderTransactionID: converter.PgTextToStringPtr(p.ProviderTransactionID),
		Status:                constant.PaymentStatus(p.PaymentStatus),
		CreatedAt:             p.CreatedAt.Time,
		UpdatedAt:             p.UpdatedAt.Time,
	}
}

func toDomainPayments(ps []sqlc.Payment) []domain.Payment {
	payments := make([]domain.Payment, len(ps))
	for i, p := range ps {
		payments[i] = *toDomain(&p)
	}
	return payments
}

// toDomain converts sqlc.PaymentRefund to domain.PaymentRefund
func toDomainRefund(r *sqlc.RefundPayment) *domain.PaymentRefund {
	return &domain.PaymentRefund{
		ID:           converter.PgUUIDToString(r.ID),
		PaymentID:    converter.PgUUIDToString(r.PaymentID),
		OrderID:      converter.PgUUIDToString(r.OrderID),
		Amount:       converter.PgNumericToFloat64(r.Amount),
		RefundStatus: constant.RefundStatus(r.RefundStatus),
		Reason:       r.Reason.String,
		CreatedAt:    r.CreatedAt.Time,
		UpdatedAt:    r.UpdatedAt.Time,
	}
}

func toDomainRefunds(rs []sqlc.RefundPayment) []domain.PaymentRefund {
	refunds := make([]domain.PaymentRefund, len(rs))
	for i, r := range rs {
		refunds[i] = *toDomainRefund(&r)
	}
	return refunds
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

func (r *paymentRepository) CreatePaymentRefund(ctx context.Context, params sqlc.CreateRefundPaymentParams) (*domain.PaymentRefund, error) {
	result, err := r.queries.CreateRefundPayment(ctx, params)
	if err != nil {
		return nil, err
	}
	return toDomainRefund(&result), nil
}

func (r *paymentRepository) GetRefundByPaymentID(ctx context.Context, paymentID string) (*domain.PaymentRefund, error) {
	result, err := r.queries.GetRefundPaymentByID(ctx, converter.StringToPgUUID(paymentID))
	if err != nil {
		return nil, err
	}
	return toDomainRefund(&result), nil
}

func (r *paymentRepository) UpdateRefundPaymentStatus(ctx context.Context, params sqlc.UpdateRefundPaymentStatusParams) (*domain.PaymentRefund, error) {
	result, err := r.queries.UpdateRefundPaymentStatus(ctx, params)
	if err != nil {
		return nil, err
	}
	return toDomainRefund(&result), nil
}

func (r *paymentRepository) GetBatchRefundPaymentsByStatus(ctx context.Context, status sqlc.RefundStatus) ([]domain.PaymentRefund, error) {
	results, err := r.queries.GetBatchRefundPaymentsByStatus(ctx, status)
	if err != nil {
		return nil, err
	}

	return toDomainRefunds(results), nil
}

func (r *paymentRepository) GetBatchPendingPayments(ctx context.Context) ([]domain.Payment, error) {
	results, err := r.queries.GetBatchPendingPayments(ctx)
	if err != nil {
		return nil, err
	}

	return toDomainPayments(results), nil
}
