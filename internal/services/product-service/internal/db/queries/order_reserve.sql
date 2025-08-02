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