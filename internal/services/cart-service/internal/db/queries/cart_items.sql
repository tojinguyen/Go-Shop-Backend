-- name: GetItemsByCartID :many
SELECT * FROM cart_items
WHERE cart_id = $1;

-- name: AddItemToCart :one
INSERT INTO cart_items (cart_id, shop_id, product_id, quantity)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteItemFromCart :exec
DELETE FROM cart_items
WHERE id = $1;

-- name: DeleteAllItemsFromCart :exec
DELETE FROM cart_items
WHERE cart_id = $1;

-- name: UpdateItemQuantity :one
UPDATE cart_items
SET quantity = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;


-- name: UpsertItemInCart :one
INSERT INTO cart_items (cart_id, shop_id, product_id, quantity)
VALUES ($1, $2, $3, $4)
ON CONFLICT (cart_id, product_id) DO UPDATE 
SET quantity = EXCLUDED.quantity, updated_at = NOW()
RETURNING *;