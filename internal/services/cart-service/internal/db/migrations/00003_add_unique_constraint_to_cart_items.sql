-- +goose Up
-- +goose StatementBegin
ALTER TABLE cart_items
ADD CONSTRAINT unique_cart_item UNIQUE (cart_id, shop_id, product_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE cart_items
DROP CONSTRAINT unique_cart_item;
-- +goose StatementEnd
