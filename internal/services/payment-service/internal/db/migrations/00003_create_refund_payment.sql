-- +goose Up
-- +goose StatementBegin
CREATE TYPE refund_status AS ENUM (
    'PENDING',
    'REFUND_REQUESTED',
    'COMPLETED',
    'FAILED'
);

CREATE TABLE refund_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id UUID NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    reason TEXT,
    provider_refund_id VARCHAR(255),
    refund_status refund_status NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (payment_id) REFERENCES payments(id)
);

ALTER TABLE payments ADD COLUMN provider_refund_id VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE payments DROP COLUMN IF EXISTS provider_refund_id;
DROP TABLE IF EXISTS refund_payments;
DROP TYPE IF EXISTS refund_status;
-- +goose StatementEnd
