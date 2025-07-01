-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS shipper_profiles (
    user_id UUID PRIMARY KEY,
    vehicle_type VARCHAR(100),
    vehicle_image_url VARCHAR(500),
    identify_card_url VARCHAR(500),
    license_plate VARCHAR(20),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING', -- Enum: PENDING, APPROVED, REJECTED

    CONSTRAINT fk_shipper_user FOREIGN KEY (user_id) REFERENCES user_profiles(user_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS shipper_profiles;
-- +goose StatementEnd
