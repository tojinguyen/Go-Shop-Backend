-- name: CreateAddress :one
INSERT INTO addresses (
  user_id,
  is_default,
  street,
  ward,
  district,
  city,
  country,
  lat,
  long,
  created_at,
  updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, now(), now()
)
RETURNING
  id,
  user_id,
  is_default,
  street,
  ward,
  district,
  city,
  country,
  lat,
  long,
  deleted_at,
  created_at,
  updated_at;

-- name: GetAddressById :one
SELECT
  id,
  user_id,
  is_default,
  street,
  ward,
  district,
  city,
  country,
  lat,
  long,
  deleted_at,
  created_at,
  updated_at
FROM addresses
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetAddressesByUserId :many
SELECT
  id,
  user_id,
  is_default,
  street,
  ward,
  district,
  city,
  country,
  lat,
  long,
  deleted_at,
  created_at,
  updated_at
FROM addresses
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY is_default DESC, created_at DESC;

-- name: UpdateAddress :one
UPDATE addresses
SET
  is_default = $2,
  street = $3,
  ward = $4,
  district = $5,
  city = $6,
  country = $7,
  lat = $8,
  long = $9,
  updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING
  id,
  user_id,
  is_default,
  street,
  ward,
  district,
  city,
  country,
  lat,
  long,
  deleted_at,
  created_at,
  updated_at;

-- name: DeleteAddress :exec
UPDATE addresses
SET
  deleted_at = now(),
  updated_at = now()
WHERE id = $1;

-- name: GetDefaultAddressByUserId :one
SELECT
  id,
  user_id,
  is_default,
  street,
  ward,
  district,
  city,
  country,
  lat,
  long,
  deleted_at,
  created_at,
  updated_at
FROM addresses
WHERE user_id = $1 AND is_default = true AND deleted_at IS NULL;

-- name: SetDefaultAddress :exec
UPDATE addresses
SET
  is_default = CASE WHEN id = $2 THEN true ELSE false END,
  updated_at = now()
WHERE user_id = $1 AND deleted_at IS NULL;
