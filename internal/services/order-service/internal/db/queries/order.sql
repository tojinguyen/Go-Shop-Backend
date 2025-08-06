-- name: CreateOrder :one 
INSERT INTO orders 
(
    id,
    owner_id,
    shop_id,
    shipping_address_id,
    promotion_id,
    shipping_fee,
    discount_amount,
    total_amount,
    final_amount,
    order_status
)
VALUES 
(
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: GetOrderByIDWithItems :one
SELECT
  o.*,
  COALESCE(
    (SELECT json_agg(oi.*)
     FROM order_items oi
     WHERE oi.order_id = o.id),
    '[]'::json
  ) as items
FROM orders o
WHERE o.id = $1;

-- name: GetOrderByID :one
SELECT * FROM orders WHERE id = $1;

-- name: GetOrdersByUserIDWithItems :many
SELECT
  o.*,
  COALESCE(
    (SELECT json_agg(oi.*)
     FROM order_items oi
     WHERE oi.order_id = o.id),
    '[]'::json
  ) as items
FROM orders o
WHERE o.owner_id = $1
ORDER BY o.created_at DESC -- Sắp xếp theo đơn hàng mới nhất
LIMIT $2 -- Giới hạn số lượng đơn hàng trên mỗi trang
OFFSET $3; -- Bỏ qua bao nhiêu đơn hàng (để qua trang mới)

-- name: UpdateOrderStatus :one
UPDATE orders
SET order_status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetStaleOrders :many
SELECT * FROM orders 
WHERE order_status = 'PENDING' OR order_status = "REFUNDED"
AND updated_at < $1
LIMIT $2;