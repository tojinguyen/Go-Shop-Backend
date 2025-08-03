-- +goose Up
-- +goose StatementBegin
CREATE TYPE outbox_event_status AS ENUM (
    'PENDING', -- Chờ được gửi
    'SENT',    -- Đã gửi thành công
    'FAILED'   -- Gửi thất bại sau nhiều lần thử, cần can thiệp thủ công
);

CREATE TABLE payment_outbox_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id UUID NOT NULL,
    event_type VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    event_status outbox_event_status NOT NULL DEFAULT 'PENDING',
    retry_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payment_outbox_events;
DROP TYPE IF EXISTS outbox_event_status;
-- +goose StatementEnd
