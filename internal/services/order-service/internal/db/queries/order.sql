-- name: CreateOrder :one 
INSERT INTO orders 
(
    user_id,
    shop_id,
    shipping_address_id,
    billing_address_id,
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