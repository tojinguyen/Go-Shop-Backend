-- name: CreateUserProfile :one
INSERT INTO user_profiles (
  user_id,
  email,
  full_name,
  birthday,
  phone,
  banned_at,
  avatar_url,
  gender,
  default_address_id,
  created_at,
  updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, now(), now()
)
RETURNING
  user_id,
  email,
  full_name,
  birthday,
  phone,
  banned_at,
  avatar_url,
  gender,
  default_address_id,
  created_at,
  updated_at;