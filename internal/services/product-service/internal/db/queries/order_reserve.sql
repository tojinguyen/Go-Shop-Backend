-- name: IsOrderReserved :one
SELECT EXISTS (
    SELECT 1 FROM order_reservations
    WHERE order_id = $1 AND reservation_status = 'RESERVED'
) AS reserved;

-- name: GetReservationStatusOfOrder :one
SELECT * FROM order_reservations
WHERE order_id = $1;

-- name: GetReservationStatusOfOrders :many
SELECT DISTINCT ON (r.order_id)
			r.order_id,
			r.shop_id,
			r.reservation_status as status,
			r.created_at,
			r.updated_at,
			true as founded
		FROM order_reservations r
		WHERE r.order_id = ANY($1::uuid[])
		ORDER BY r.order_id, r.created_at DESC;

-- name: UpdateReservationStatusOfOrder :exec
UPDATE order_reservations
SET reservation_status = $2, updated_at = NOW()
WHERE order_id = $1;

-- name: ReserveOrder :exec
INSERT INTO order_reservations (
	order_id,
	shop_id,
	reservation_status
)
VALUES ($1, $2, $3);