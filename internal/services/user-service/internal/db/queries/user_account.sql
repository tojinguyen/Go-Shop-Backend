-- name: CreateUserAccount :one
INSERT INTO user_accounts (
  email, hashed_password
) VALUES (
  $1, $2
) RETURNING id, email, created_at, updated_at;