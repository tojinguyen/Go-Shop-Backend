-- +goose Up
-- +goose StatementBegin
CREATE TYPE reservation_status AS ENUM (
    'RESERVED',   -- Hàng đã được đặt trước, đang chờ thanh toán
    'COMMITTED',  -- Đơn hàng đã thanh toán, số lượng đã được trừ khỏi kho chính
    'CANCELLED'   -- Đơn hàng bị hủy, hàng đã được trả lại kho
);

CREATE TABLE product_reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL, -- Sẽ có nhiều dòng cho cùng một order_id
    product_id UUID NOT NULL,
    shop_id UUID NOT NULL, -- Vẫn hữu ích để truy vấn theo shop
    quantity_reserved INT NOT NULL CHECK (quantity_reserved > 0),
    status reservation_status NOT NULL DEFAULT 'RESERVED',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Ràng buộc quan trọng nhất:
    -- Đảm bảo cho một đơn hàng, một sản phẩm chỉ được đặt trước MỘT LẦN.
    -- Đây chính là chìa khóa cho tính Idempotency.
    CONSTRAINT uq_order_product UNIQUE (order_id, product_id),

    -- Foreign key đến bảng products
    CONSTRAINT fk_reservation_product FOREIGN KEY (product_id) REFERENCES products(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE reservation_status;
-- +goose StatementEnd
