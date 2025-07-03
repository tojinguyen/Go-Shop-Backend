-- name: CreateUserProfile :one
INSERT INTO user_profiles (
  user_id,
  email,
  full_name,
  birthday,
  phone,
  user_role,
  banned_at,
  avatar_url,
  gender,
  default_address_id,
  created_at,
  updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, now(), now()
)
RETURNING
  user_id,
  email,
  full_name,
  birthday,
  phone,
  user_role,
  banned_at,
  avatar_url,
  gender,
  default_address_id,
  created_at,
  updated_at;

-- name: GetUserProfileByUserId :one
SELECT
  user_id,
  email,
  full_name,
  birthday,
  phone,
  user_role,
  banned_at,
  avatar_url,
  gender,
  default_address_id,
  created_at,
  updated_at
FROM user_profiles
WHERE user_id = $1 AND banned_at IS NULL;

-- name: UpdateUserProfile :one
UPDATE user_profiles
SET
  email = $2,
  full_name = $3,
  birthday = $4,
  phone = $5,
  user_role = $6,
  banned_at = $7,
  avatar_url = $8,
  gender = $9,
  default_address_id = $10,
  updated_at = now()
WHERE user_id = $1
RETURNING
  user_id,
  email,
  full_name,
  birthday,
  phone,
  user_role,
  banned_at,
  avatar_url,
  gender,
  default_address_id,
  created_at,
  updated_at;


-- name: SoftDeleteUserProfile :exec
UPDATE user_profiles
SET
  banned_at = now(),
  updated_at = now()
WHERE user_id = $1;