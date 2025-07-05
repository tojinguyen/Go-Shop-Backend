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
  created_at,
  updated_at
FROM shops
WHERE id = $1;