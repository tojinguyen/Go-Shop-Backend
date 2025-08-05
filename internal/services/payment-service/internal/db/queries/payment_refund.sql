-- name: CreateRefundPayment :one
INSERT INTO refund_payments (
    payment_id,
    order_id,
    amount,
    reason,
    provider_refund_id,
    refund_status
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetRefundPaymentByID :one
SELECT * FROM refund_payments WHERE id = $1;

-- name: UpdateRefundPaymentStatus :exec
UPDATE refund_payments SET refund_status = $2 WHERE id = $1;