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
WHERE id = $1 AND deleted_at IS NULL;

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