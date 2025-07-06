-- name: CreatePromotion :one
INSERT INTO shop_promotions (
    id,
    shop_id,
    promotion_name,
    promotion_type,
    discount_value,
    max_discount_amount,
    min_purchase_amount,
    usage_limit_per_user,
    start_time,
    end_time,
    promotion_status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetPromotionByID :one
SELECT
    id,
    shop_id,
    promotion_name,
    promotion_type,
    discount_value,
    max_discount_amount,
    min_purchase_amount,
    usage_limit_per_user,
    start_time,
    end_time,
    promotion_status,
    created_at,
    updated_at
FROM shop_promotions
WHERE id = $1;

-- name: GetPromotionsByShopID :many
SELECT
    id,
    shop_id,
    promotion_name,
    promotion_type,
    discount_value,
    max_discount_amount,
    min_purchase_amount,
    usage_limit_per_user,
    start_time,
    end_time,
    promotion_status,
    created_at,
    updated_at
FROM shop_promotions
WHERE shop_id = $1;

-- name: GetAllPromotionsByStatus :many
SELECT 
    id,
    shop_id,
    promotion_name,
    promotion_type,
    discount_value,
    max_discount_amount,
    min_purchase_amount,
    usage_limit_per_user,
    start_time,
    end_time,
    promotion_status,
    created_at,
    updated_at
FROM shop_promotions
WHERE promotion_status = $1;

-- name: UpdatePromotion :one
UPDATE shop_promotions
SET
    promotion_name = $2,
    promotion_type = $3,
    discount_value = $4,
    max_discount_amount = $5,
    min_purchase_amount = $6,
    usage_limit_per_user = $7,
    start_time = $8,
    end_time = $9,
    promotion_status = $10,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeletePromotion :exec
DELETE FROM shop_promotions
WHERE id = $1;
