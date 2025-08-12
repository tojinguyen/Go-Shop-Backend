-- name: CreateShopAddress :one
INSERT INTO addresses (
  id,
  shop_id,
  street,
  ward,
  district,
  city,
  country,
  lat,
  long
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING id;