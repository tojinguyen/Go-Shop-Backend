-- name: CreateShipper :one
INSERT INTO shipper_profiles (
    user_id,
    vehicle_type,
    vehicle_image_url,
    identify_card_url,
    license_plate
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetShipperByUserID :one
SELECT * FROM shipper_profiles
WHERE user_id = $1;

-- name: UpdateShipperByUserID :one
UPDATE shipper_profiles 
SET 
    vehicle_type = $2,
    vehicle_image_url = $3,
    identify_card_url = $4,
    license_plate = $5
WHERE user_id = $1
RETURNING *;


-- name: DeleteShipperByUserID :exec
DELETE FROM shipper_profiles
WHERE user_id = $1;
