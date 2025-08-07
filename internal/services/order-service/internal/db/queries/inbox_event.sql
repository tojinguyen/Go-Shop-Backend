-- name: CreateInboxEvent :one
INSERT INTO order_inbox_events (
    event_id,
    event_type,
    source_service,
    payload,
    event_status
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetPendingInboxEvents :many
SELECT * FROM order_inbox_events
WHERE event_status = 'PENDING'
ORDER BY received_at ASC
LIMIT $1;

-- name: GetInboxEventByEventId :one
SELECT * FROM order_inbox_events
WHERE event_id = $1;

-- name: UpdateInboxEventStatus :one
UPDATE order_inbox_events
SET
    event_status = $2,
    retry_count = $3,
    processed_at = CASE 
        WHEN $2 = 'PROCESSED' THEN NOW() 
        ELSE processed_at 
    END,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetFailedInboxEvents :many
SELECT * FROM order_inbox_events
WHERE event_status = 'FAILED' AND retry_count < max_retry
ORDER BY received_at ASC
LIMIT $1;

-- name: GetInboxEventStats :one
SELECT 
    COUNT(*) FILTER (WHERE event_status = 'PENDING') as pending_count,
    COUNT(*) FILTER (WHERE event_status = 'PROCESSED') as processed_count,
    COUNT(*) FILTER (WHERE event_status = 'FAILED') as failed_count,
    COUNT(*) as total_count
FROM order_inbox_events;

-- name: CleanupOldInboxEvents :exec
DELETE FROM order_inbox_events
WHERE event_status = 'PROCESSED' 
AND processed_at < NOW() - INTERVAL '30 days';
