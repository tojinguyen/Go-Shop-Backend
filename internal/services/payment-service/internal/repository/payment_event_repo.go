package repository

import (
	"context"

	"github.com/toji-dev/go-shop/internal/pkg/converter"
	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/domain"
)

type PaymentEventRepository interface {
	CreatePaymentEvent(ctx context.Context, paymentEvent *domain.PaymentEvent) (*domain.PaymentEvent, error)
	GetBatchPaymentEventByEventTypeAndStatus(ctx context.Context, eventType domain.PaymentEventType, status domain.PaymentEventStatus, limit int) ([]*domain.PaymentEvent, error)
	UpdatePaymentEvent(ctx context.Context, paymentEvent *domain.PaymentEvent) (*domain.PaymentEvent, error)
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

func (r *paymentEventRepository) CreatePaymentEvent(ctx context.Context, paymentEvent *domain.PaymentEvent) (*domain.PaymentEvent, error) {
	params := sqlc.CreatePaymentEventParams{
		PaymentID:   converter.StringToPgUUID(paymentEvent.PaymentID),
		OrderID:     converter.StringToPgUUID(paymentEvent.OrderID),
		EventType:   paymentEvent.EventType,
		Payload:     []byte(paymentEvent.Payload),
		EventStatus: sqlc.OutboxEventStatus(sqlc.PaymentStatusPENDING),
	}

	result, err := r.queries.CreatePaymentEvent(ctx, params)
	if err != nil {
		return nil, err
	}

	eventResult := ToDomain(&result)

	return eventResult, nil
}

func (r *paymentEventRepository) GetBatchPaymentEventByEventTypeAndStatus(ctx context.Context, eventType domain.PaymentEventType, status domain.PaymentEventStatus, limit int) ([]*domain.PaymentEvent, error) {
	params := sqlc.GetBatchPaymentEventsByEventTypeAndStatusParams{
		EventType:   string(eventType),
		EventStatus: sqlc.OutboxEventStatus(status),
		Limit:       int32(limit),
	}

	result, err := r.queries.GetBatchPaymentEventsByEventTypeAndStatus(ctx, params)
	if err != nil {
		return nil, err
	}

	paymentEvents := make([]*domain.PaymentEvent, 0, len(result))
	for _, item := range result {
		paymentEvents = append(paymentEvents, ToDomain(&item))
	}

	return paymentEvents, nil
}

func (r *paymentEventRepository) UpdatePaymentEvent(ctx context.Context, paymentEvent *domain.PaymentEvent) (*domain.PaymentEvent, error) {
	params := sqlc.UpdatePaymentEventParams{
		ID:          converter.StringToPgUUID(paymentEvent.ID),
		EventStatus: sqlc.OutboxEventStatus(paymentEvent.EventStatus),
		RetryCount:  int32(paymentEvent.RetryCount),
	}

	result, err := r.queries.UpdatePaymentEvent(ctx, params)
	if err != nil {
		return nil, err
	}

	eventResult := ToDomain(&result)

	return eventResult, nil
}

func ToDomain(pgPaymentEvent *sqlc.PaymentOutboxEvent) *domain.PaymentEvent {
	if pgPaymentEvent == nil {
		return nil
	}

	return &domain.PaymentEvent{
		ID:          pgPaymentEvent.ID.String(),
		PaymentID:   pgPaymentEvent.PaymentID.String(),
		OrderID:     pgPaymentEvent.OrderID.String(),
		EventType:   pgPaymentEvent.EventType,
		Payload:     string(pgPaymentEvent.Payload),
		EventStatus: domain.PaymentEventStatus(pgPaymentEvent.EventStatus),
		RetryCount:  int(pgPaymentEvent.RetryCount),
		CreatedAt:   *converter.PgTimeToTimePtr(pgPaymentEvent.CreatedAt),
		UpdatedAt:   *converter.PgTimeToTimePtr(pgPaymentEvent.UpdatedAt),
	}
}
