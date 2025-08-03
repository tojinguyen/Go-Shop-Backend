-- name: CreatePaymentEvent :one
INSERT INTO payment_outbox_events (
    payment_id,
    order_id,
    event_type,
    payload,
    event_status
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetBatchPaymentEventsByEventTypeAndStatus :many
SELECT * FROM payment_outbox_events
WHERE event_type = $1 AND event_status = $2
ORDER BY created_at ASC
LIMIT $3;

-- name: UpdatePaymentEvent :one
UPDATE payment_outbox_events
SET
    event_status = $2,
    payload = $3,
    retry_count = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdatePaymentEventStatus :one
UPDATE payment_outbox_events
SET
    event_status = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
