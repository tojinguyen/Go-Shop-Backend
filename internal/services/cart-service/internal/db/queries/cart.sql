-- name: GetCartByOwnerID :one
SELECT * FROM carts
WHERE owner_id = $1;

-- name: CreateCart :one
INSERT INTO carts (owner_id)
VALUES ($1)
RETURNING *;