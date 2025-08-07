package repository

import (
	"context"

	"github.com/toji-dev/go-shop/internal/pkg/converter"
	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/order-service/internal/domain"
)

type InboxEventRepository interface {
	CreateInboxEvent(ctx context.Context, event *domain.InboxEvent) (*domain.InboxEvent, error)
	GetPendingInboxEvents(ctx context.Context, limit int) ([]*domain.InboxEvent, error)
	GetInboxEventByEventID(ctx context.Context, eventID string) (*domain.InboxEvent, error)
	UpdateInboxEventStatus(ctx context.Context, event *domain.InboxEvent) (*domain.InboxEvent, error)
	GetFailedInboxEvents(ctx context.Context, limit int) ([]*domain.InboxEvent, error)
	GetInboxEventStats(ctx context.Context) (*InboxEventStats, error)
	CleanupOldInboxEvents(ctx context.Context) error
}

type InboxEventStats struct {
	PendingCount   int64 `json:"pending_count"`
	ProcessedCount int64 `json:"processed_count"`
	FailedCount    int64 `json:"failed_count"`
	TotalCount     int64 `json:"total_count"`
}

type inboxEventRepository struct {
	db      *postgresql_infra.PostgreSQLService
	queries *sqlc.Queries
}

func NewInboxEventRepository(db *postgresql_infra.PostgreSQLService) InboxEventRepository {
	if db == nil {
		return nil
	}

	queries := sqlc.New(db.GetPool())

	return &inboxEventRepository{
		db:      db,
		queries: queries,
	}
}

func (r *inboxEventRepository) CreateInboxEvent(ctx context.Context, event *domain.InboxEvent) (*domain.InboxEvent, error) {
	params := sqlc.CreateInboxEventParams{
		EventID:       event.EventID,
		EventType:     event.EventType,
		SourceService: event.SourceService,
		Payload:       []byte(event.Payload),
		EventStatus:   sqlc.InboxEventStatus(event.EventStatus),
	}

	result, err := r.queries.CreateInboxEvent(ctx, params)
	if err != nil {
		return nil, err
	}

	return inboxEventToDomain(&result), nil
}

func (r *inboxEventRepository) GetPendingInboxEvents(ctx context.Context, limit int) ([]*domain.InboxEvent, error) {
	results, err := r.queries.GetPendingInboxEvents(ctx, int32(limit))
	if err != nil {
		return nil, err
	}

	events := make([]*domain.InboxEvent, 0, len(results))
	for _, result := range results {
		events = append(events, inboxEventToDomain(&result))
	}

	return events, nil
}

func (r *inboxEventRepository) GetInboxEventByEventID(ctx context.Context, eventID string) (*domain.InboxEvent, error) {
	result, err := r.queries.GetInboxEventByEventId(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return inboxEventToDomain(&result), nil
}

func (r *inboxEventRepository) UpdateInboxEventStatus(ctx context.Context, event *domain.InboxEvent) (*domain.InboxEvent, error) {
	params := sqlc.UpdateInboxEventStatusParams{
		ID:          converter.StringToPgUUID(event.ID),
		EventStatus: sqlc.InboxEventStatus(event.EventStatus),
		RetryCount:  int32(event.RetryCount),
	}

	result, err := r.queries.UpdateInboxEventStatus(ctx, params)
	if err != nil {
		return nil, err
	}

	return inboxEventToDomain(&result), nil
}

func (r *inboxEventRepository) GetFailedInboxEvents(ctx context.Context, limit int) ([]*domain.InboxEvent, error) {
	results, err := r.queries.GetFailedInboxEvents(ctx, int32(limit))
	if err != nil {
		return nil, err
	}

	events := make([]*domain.InboxEvent, 0, len(results))
	for _, result := range results {
		events = append(events, inboxEventToDomain(&result))
	}

	return events, nil
}

func (r *inboxEventRepository) GetInboxEventStats(ctx context.Context) (*InboxEventStats, error) {
	result, err := r.queries.GetInboxEventStats(ctx)
	if err != nil {
		return nil, err
	}

	return &InboxEventStats{
		PendingCount:   result.PendingCount,
		ProcessedCount: result.ProcessedCount,
		FailedCount:    result.FailedCount,
		TotalCount:     result.TotalCount,
	}, nil
}

func (r *inboxEventRepository) CleanupOldInboxEvents(ctx context.Context) error {
	return r.queries.CleanupOldInboxEvents(ctx)
}

// Helper function to convert SQLC model to domain model
func inboxEventToDomain(sqlcEvent *sqlc.OrderInboxEvent) *domain.InboxEvent {
	if sqlcEvent == nil {
		return nil
	}

	return &domain.InboxEvent{
		ID:            sqlcEvent.ID.String(),
		EventID:       sqlcEvent.EventID,
		EventType:     sqlcEvent.EventType,
		SourceService: sqlcEvent.SourceService,
		Payload:       string(sqlcEvent.Payload),
		EventStatus:   domain.InboxEventStatus(sqlcEvent.EventStatus),
		RetryCount:    int(sqlcEvent.RetryCount),
		MaxRetry:      int(sqlcEvent.MaxRetry),
		ReceivedAt:    *converter.PgTimeToTimePtr(sqlcEvent.ReceivedAt),
		ProcessedAt:   converter.PgTimeToTimePtr(sqlcEvent.ProcessedAt),
		CreatedAt:     *converter.PgTimeToTimePtr(sqlcEvent.CreatedAt),
		UpdatedAt:     *converter.PgTimeToTimePtr(sqlcEvent.UpdatedAt),
	}
}
