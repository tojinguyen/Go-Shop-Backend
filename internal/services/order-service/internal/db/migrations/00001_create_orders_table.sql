-- +goose Up
-- +goose StatementBegin
CREATE TYPE order_status AS ENUM 
(
    'PENDING',
    'PENDING_PAYMENT',
    'PAYMENT_FAILED',
    'PROCESSING',
    'SHIPPED',
    'DELIVERING',
    'DELIVERED',
    'CANCELED',
    'FAILED'
);

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL,
    shop_id UUID NOT NULL,
    shipping_address_id UUID NOT NULL,
    promotion_id UUID,

    shipping_fee NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    discount_amount NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    total_amount NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    final_amount NUMERIC(10, 2) NOT NULL DEFAULT 0.00,

    order_status order_status NOT NULL DEFAULT 'PENDING',
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
DROP TYPE IF EXISTS order_status;
-- +goose StatementEnd
