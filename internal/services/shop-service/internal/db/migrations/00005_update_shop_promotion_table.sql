-- +goose Up
-- +goose StatementBegin
ALTER TABLE shop_promotions
ADD COLUMN usage_limit_total INTEGER CHECK (usage_limit_total > 0),
ADD COLUMN current_usage_count INTEGER DEFAULT 0 CHECK (current_usage_count >= 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE shop_promotions
DROP COLUMN current_usage_count,
DROP COLUMN usage_limit_total;
-- +goose StatementEnd
