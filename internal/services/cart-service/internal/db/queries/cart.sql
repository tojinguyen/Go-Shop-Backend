-- name: GetCartByOwnerID :one
SELECT * FROM carts
WHERE owner_id = $1;

-- name: CreateCart :one
INSERT INTO carts (owner_id)
VALUES ($1)
RETURNING *;

-- name: UpsertCart :one
INSERT INTO carts (owner_id, updated_at)
VALUES ($1, NOW())
ON CONFLICT (owner_id) DO UPDATE 
SET updated_at = NOW()
RETURNING *;

-- name: DeleteCart :exec
DELETE FROM carts
WHERE owner_id = $1;