-- +goose Up
-- +goose StatementBegin
CREATE TYPE order_status AS ENUM 
(
    'PENDING_PAYMENT',
    'PAYMENT_FAILED',
    'PROCESSING',
    'SHIPPED',
    'DELIVERING',
    'DELIVERED',
    'CANCELED'
);

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    shop_id UUID NOT NULL,
    shipping_address_id UUID NOT NULL,
    billing_address_id UUID NOT NULL,
    promotion_id UUID,

    order_status order_status NOT NULL DEFAULT 'PENDING_PAYMENT',
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
DROP TYPE IF EXISTS order_status;
-- +goose StatementEnd
