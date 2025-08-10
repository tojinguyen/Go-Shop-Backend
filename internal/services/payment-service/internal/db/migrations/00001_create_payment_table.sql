-- +goose Up
-- +goose StatementBegin
CREATE TYPE payment_status AS ENUM (
    'PENDING',
    'PROCESSING',
    'SUCCESS',
    'FAILED',
    'REFUNDED'
);

CREATE TYPE payment_method AS ENUM (
    'COD',
    'CREDIT_CARD',
    'BANK_TRANSFER',
    'E_WALLET'
);

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL UNIQUE,
    user_id UUID NOT NULL,
    amount NUMERIC(12, 2) NOT NULL CHECK (amount >= 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    payment_method payment_method NOT NULL,
    payment_provider VARCHAR(50), -- e.g., 'STRIPE', 'MOMO', 'VNPAY', 'NONE'
    provider_transaction_id VARCHAR(255),
    payment_status payment_status NOT NULL DEFAULT 'PENDING',
    request_id VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payments;
DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS payment_method;
-- +goose StatementEnd
