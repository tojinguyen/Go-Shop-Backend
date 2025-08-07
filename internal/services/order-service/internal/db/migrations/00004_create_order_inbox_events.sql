-- +goose Up
-- +goose StatementBegin

-- Enum cho trạng thái inbox event
CREATE TYPE inbox_event_status AS ENUM (
    'PENDING',
    'PROCESSED',
    'FAILED'
);

-- Bảng inbox events cho Order Service
CREATE TABLE order_inbox_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Thông tin về event
    event_id VARCHAR(255) UNIQUE NOT NULL, -- UUID từ hệ thống gửi (để tránh duplicate)
    event_type VARCHAR(100) NOT NULL,      -- Loại event: PAYMENT_SUCCESS, REFUND_SUCCEEDED, etc.
    source_service VARCHAR(50) NOT NULL,   -- Service gửi event: payment-service, etc.
    
    -- Data của event
    payload JSONB NOT NULL,                -- Dữ liệu event dạng JSON
    
    -- Processing information
    event_status inbox_event_status NOT NULL DEFAULT 'PENDING',
    retry_count INTEGER NOT NULL DEFAULT 0,
    max_retry INTEGER NOT NULL DEFAULT 5,
    
    -- Thông tin thời gian
    received_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes để tối ưu query
CREATE INDEX idx_order_inbox_events_status ON order_inbox_events(event_status);
CREATE INDEX idx_order_inbox_events_type_status ON order_inbox_events(event_type, event_status);
CREATE INDEX idx_order_inbox_events_received_at ON order_inbox_events(received_at);
CREATE INDEX idx_order_inbox_events_retry ON order_inbox_events(retry_count) WHERE event_status = 'FAILED';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_inbox_events;
DROP TYPE IF EXISTS inbox_event_status;
-- +goose StatementEnd
