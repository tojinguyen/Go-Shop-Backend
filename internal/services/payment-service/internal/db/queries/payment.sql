-- name: CreatePayment :one
INSERT INTO payments (
    order_id,
    user_id,
    amount,
    currency,
    payment_method,
    payment_provider,
    payment_status,
    request_id
) VALUES (
    $1, $2, $3, $4, $5, $6, 'PENDING', $7
) RETURNING *;

-- name: UpdatePaymentStatus :one
UPDATE payments
SET
    payment_status = $2,
    provider_transaction_id = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetPaymentByOrderID :one
SELECT * FROM payments
WHERE order_id = $1;

-- name: UpdatePaymentProviderRefundID :one
UPDATE payments
SET
    provider_refund_id = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetBatchPendingPayments :many
SELECT * FROM payments
WHERE payment_status = 'PENDING' AND created_at < NOW() - INTERVAL '15 minutes';