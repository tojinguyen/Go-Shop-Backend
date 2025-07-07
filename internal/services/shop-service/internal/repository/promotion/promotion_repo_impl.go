package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/toji-dev/go-shop/internal/pkg/converter"
	postgresql_infra "github.com/toji-dev/go-shop/internal/pkg/infra/postgreql-infra"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/db/sqlc"
	"github.com/toji-dev/go-shop/internal/services/shop-service/internal/domain"
)

type PostgresPromotionRepository struct {
	db      *postgresql_infra.PostgreSQLService
	queries *sqlc.Queries
}

func NewPostgresPromotionRepository(db *postgresql_infra.PostgreSQLService) PromotionRepository {
	return &PostgresPromotionRepository{
		db:      db,
		queries: sqlc.New(db.GetPool()),
	}
}

func (r *PostgresPromotionRepository) Create(ctx context.Context, p *domain.Promotion) error {
	params := sqlc.CreatePromotionParams{
		ID:                converter.UUIDToPgUUID(p.ID),
		ShopID:            converter.UUIDToPgUUID(p.ShopID),
		PromotionName:     p.PromotionName,
		PromotionType:     sqlc.PromotionType(p.PromotionType),
		DiscountValue:     converter.Float64ToPgNumeric(p.DiscountValue),
		MaxDiscountAmount: converter.Float64ToPgNumeric(*p.MaxDiscountAmount),
		MinPurchaseAmount: converter.Float64ToPgNumeric(p.MinPurchaseAmount),
		UsageLimitPerUser: converter.Int32ToPgInt4(p.UsageLimitPerUser),
		StartTime:         converter.TimeToPgTime(p.StartTime),
		EndTime:           converter.TimeToPgTime(p.EndTime),
		PromotionStatus:   sqlc.NullPromotionStatus{PromotionStatus: sqlc.PromotionStatus(p.PromotionStatus), Valid: true},
	}
	_, err := r.queries.CreatePromotion(ctx, params)
	if err != nil {
		log.Printf("Error creating promotion: %v", err)
		return fmt.Errorf("failed to create promotion: %w", err)
	}
	return nil
}

func (r *PostgresPromotionRepository) GetByID(ctx context.Context, id string) (*domain.Promotion, error) {
	promo, err := r.queries.GetPromotionByID(ctx, converter.StringToPgUUID(id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Return nil, nil if not found
		}
		return nil, fmt.Errorf("failed to get promotion by id: %w", err)
	}
	return r.toDomain(&promo), nil
}

func (r *PostgresPromotionRepository) GetByShopID(ctx context.Context, shopID string) ([]*domain.Promotion, error) {
	promos, err := r.queries.GetPromotionsByShopID(ctx, converter.StringToPgUUID(shopID))
	if err != nil {
		return nil, fmt.Errorf("failed to get promotions by shop id: %w", err)
	}

	domainPromos := make([]*domain.Promotion, len(promos))
	for i, promo := range promos {
		domainPromos[i] = r.toDomain(&promo)
	}
	return domainPromos, nil
}

func (r *PostgresPromotionRepository) Update(ctx context.Context, p *domain.Promotion) error {
	params := sqlc.UpdatePromotionParams{
		ID:                converter.UUIDToPgUUID(p.ID),
		PromotionName:     p.PromotionName,
		PromotionType:     sqlc.PromotionType(p.PromotionType),
		DiscountValue:     converter.Float64ToPgNumeric(p.DiscountValue),
		MaxDiscountAmount: converter.Float64ToPgNumeric(*p.MaxDiscountAmount),
		MinPurchaseAmount: converter.Float64ToPgNumeric(p.MinPurchaseAmount),
		UsageLimitPerUser: converter.Int32ToPgInt4(p.UsageLimitPerUser),
		StartTime:         converter.TimeToPgTime(p.StartTime),
		EndTime:           converter.TimeToPgTime(p.EndTime),
		PromotionStatus:   sqlc.NullPromotionStatus{PromotionStatus: sqlc.PromotionStatus(p.PromotionStatus), Valid: true},
	}
	_, err := r.queries.UpdatePromotion(ctx, params)
	return err
}

func (r *PostgresPromotionRepository) Delete(ctx context.Context, id string) error {
	return r.queries.DeletePromotion(ctx, converter.StringToPgUUID(id))
}

// toDomain is a helper to convert sqlc model to domain model
func (r *PostgresPromotionRepository) toDomain(p *sqlc.ShopPromotion) *domain.Promotion {
	return &domain.Promotion{
		ID:                converter.PgUUIDToUUID(p.ID),
		ShopID:            converter.PgUUIDToUUID(p.ShopID),
		PromotionName:     p.PromotionName,
		PromotionType:     domain.PromotionType(p.PromotionType),
		DiscountValue:     *converter.PgNumericToFloat64Ptr(p.DiscountValue),
		MaxDiscountAmount: converter.PgNumericToFloat64Ptr(p.MaxDiscountAmount),
		MinPurchaseAmount: *converter.PgNumericToFloat64Ptr(p.MinPurchaseAmount),
		UsageLimitPerUser: *converter.PgInt4ToInt32Ptr(p.UsageLimitPerUser),
		StartTime:         p.StartTime.Time,
		EndTime:           p.EndTime.Time,
		PromotionStatus:   domain.PromotionStatus(p.PromotionStatus.PromotionStatus),
		CreatedAt:         p.CreatedAt.Time,
		UpdatedAt:         p.UpdatedAt.Time,
	}
}
