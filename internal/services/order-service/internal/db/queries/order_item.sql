-- name: CreateOrderItem :one
INSERT INTO order_items (
    order_id,
    product_id,
    shop_id,
    quantity,
    price
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;