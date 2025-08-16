-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_shop_id ON products (shop_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_shop_id;
-- +goose StatementEnd
