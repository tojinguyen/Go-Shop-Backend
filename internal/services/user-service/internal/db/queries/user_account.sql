-- name: CreateUserAccount :one
INSERT INTO user_accounts (
  email, hashed_password
) VALUES (
  $1, $2
) RETURNING id, email, created_at, updated_at;

-- name: GetUserAccountByEmail :one
SELECT id, email, hashed_password, last_login_at, created_at, updated_at 
FROM user_accounts 
WHERE email = $1 AND deleted_at IS NULL;

-- name: GetUserAccountByID :one
SELECT id, email, hashed_password, last_login_at, created_at, updated_at 
FROM user_accounts 
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateLastLoginAt :exec
UPDATE user_accounts 
SET last_login_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateUserPassword :exec
UPDATE user_accounts 
SET hashed_password = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL;

-- name: SoftDeleteUserAccount :exec
UPDATE user_accounts 
SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL;

-- name: CheckUserExistsByEmail :one
SELECT id 
FROM user_accounts 
WHERE email = $1 AND deleted_at IS NULL;