-- +goose Up
-- +goose StatementBegin
CREATE TYPE reservation_status AS ENUM (
    'RESERVED',   -- Hàng đã được đặt trước, đang chờ thanh toán
    'COMMITTED',  -- Đơn hàng đã thanh toán, số lượng đã được trừ khỏi kho chính
    'CANCELLED'   -- Đơn hàng bị hủy, hàng đã được trả lại kho
);

CREATE TABLE order_reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL UNIQUE, 
    shop_id UUID NOT NULL, 
    reservation_status reservation_status NOT NULL DEFAULT 'RESERVED',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_reservations;
DROP TABLE IF EXISTS product_reservations;

DROP TYPE reservation_status;
-- +goose StatementEnd
