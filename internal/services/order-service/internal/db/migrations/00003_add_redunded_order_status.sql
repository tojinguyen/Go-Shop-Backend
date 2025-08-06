-- +goose Up
-- +goose StatementBegin
ALTER TYPE order_status ADD VALUE 'REFUNDED';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
