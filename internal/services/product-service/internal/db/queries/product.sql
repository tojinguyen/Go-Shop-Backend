-- name: CreateProduct :one
INSERT INTO products (
    shop_id,
    product_name,
    thumbnail_url,
    product_description,
    category_id,
    price,
    currency,
    quantity,
    reserve_quantity,
    product_status
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: GetProductByID :one
SELECT * FROM products
WHERE id = $1 AND delete_at IS NULL;

-- name: GetListProductsByShop :many
SELECT * FROM products
WHERE shop_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountProductsByShop :one
SELECT count(*) FROM products
WHERE shop_id = $1 AND deleted_at IS NULL;

-- name: UpdateProduct :one
UPDATE products
SET
  product_name = $2,
  product_description = $3,
  category_id = $4,
  price = $5,
  currency = $6,
  quantity = $7,
  thumbnail_url = $8,
  product_status = $9,
  updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteProduct :exec
UPDATE products
SET
  deleted_at = NOW(),
  product_status = 'DISCONTINUED', -- Hoặc một trạng thái xóa khác
  updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetProductsByIDs :many
SELECT * FROM products
WHERE id = ANY(@product_ids::uuid[]) AND delete_at IS NULL;


-- name: UpdateProductStock :one
UPDATE products
SET
    quantity = $2,
    reserve_quantity = $3,
    updated_at = NOW()
WHERE id = $1 AND delete_at IS NULL
RETURNING *;

-- name: GetProductsByIDsForUpdate :many
SELECT * FROM products
WHERE id = ANY(@product_ids::uuid[]) AND delete_at IS NULL
FOR UPDATE;

-- name: UnreserveProducts :exec
UPDATE products
SET
    reserve_quantity = CASE 
        WHEN p.reserve_quantity >= p.quantity THEN p.reserve_quantity - p.quantity
        ELSE p.reserve_quantity
    END
FROM (
    SELECT 
        CAST(unnest(@product_ids::uuid[]) as uuid) as id,
        CAST(unnest(@quantities::int[]) as integer) as quantity
) AS p
WHERE products.id = p.id AND products.delete_at IS NULL;