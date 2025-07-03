-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_promotion_usages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    promotion_id UUID NOT NULL,
    user_id UUID NOT NULL,
    used_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    order_id UUID NOT NULL,
    usage_count INTEGER DEFAULT 1 CHECK (usage_count > 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE shop_promotion_usages ADD CONSTRAINT fk_shop_promotion_usages_promotion_id 
    FOREIGN KEY (promotion_id) REFERENCES shop_promotions(id) ON DELETE CASCADE;

ALTER TABLE shop_promotion_usages ADD CONSTRAINT uk_shop_promotion_usages_order_promotion 
    UNIQUE (order_id, promotion_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS shop_promotion_usages;
-- +goose StatementEnd
