-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_profiles (
    user_id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    full_name VARCHAR(255),
    birthday DATE,
    phone VARCHAR(20) UNIQUE,
    banned_at TIMESTAMP WITH TIME ZONE,
    avatar_url TEXT,
    gender VARCHAR(10), -- Enum: MALE, FEMALE, OTHER
    default_address_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),

    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES user_accounts(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_profiles;
-- +goose StatementEnd
