-- name: CreateShop :one
INSERT INTO shops (
  id,
  owner_id,
  shop_name,
  avatar_url,
  banner_url,
  shop_description,
  address_id,
  phone,
  email,
  rating,
  active_at,
  created_at,
  updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
) RETURNING *;

-- name: GetShopByID :one
SELECT 
  id,
  owner_id,
  shop_name,
  avatar_url,
  banner_url,
  shop_description,
  address_id,
  phone,
  email,
  rating,
  active_at,
  banned_at,
  created_at,
  updated_at
FROM shops
WHERE id = $1;

-- name: GetShopsByOwnerID :many
SELECT 
  id,
  owner_id,
  shop_name,
  avatar_url,
  banner_url,
  shop_description,
  address_id,
  phone,
  email,
  rating,
  active_at,
  banned_at,
  created_at,
  updated_at
FROM shops
WHERE owner_id = $1
ORDER BY created_at DESC;

-- name: UpdateShop :one
UPDATE shops
SET 
  shop_name = $2,
  avatar_url = $3,
  banner_url = $4,
  shop_description = $5,
  address_id = $6,
  phone = $7,
  email = $8,
  rating = $9,
  active_at = $10,
  banned_at = $11,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteShop :exec
DELETE FROM shops WHERE id = $1;