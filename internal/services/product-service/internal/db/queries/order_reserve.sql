-- name: IsOrderReserved :one
SELECT EXISTS (
    SELECT 1 FROM order_reservations
    WHERE order_id = $1 AND reservation_status = 'RESERVED'
) AS reserved;