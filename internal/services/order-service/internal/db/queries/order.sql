-- name: CreateOrder :one 
INSERT INTO orders 
(
    id,
    user_id,
    shop_id,
    shipping_address_id,
    promotion_id,
    order_status
)
VALUES 
(
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1;

-- name: GetOrdersByUserID :many
SELECT * FROM orders
WHERE user_id = $1;

-- name: UpdateOrderStatus :one
UPDATE orders
SET order_status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;


-- name: GetStaleOrders :many
SELECT * FROM orders 
WHERE order_status = 'PENDING'
AND updated_at < $1
LIMIT $2;