-- +goose Up
-- +goose StatementBegin

-- 1. Tạo ENUM cho product status
CREATE TYPE product_status AS ENUM (
    'ACTIVE',
    'INACTIVE',
    'OUT_OF_STOCK',
    'DISCONTINUED',
    'BANNED'
);

-- 2. Tạo bảng products
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    shop_id UUID NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    thumbnail_url TEXT,
    product_description TEXT,
    category_id UUID,
    price NUMERIC(12, 2) NOT NULL CHECK (price >= 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    quantity INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    reserve_quantity INTEGER NOT NULL DEFAULT 0 CHECK (reserve_quantity >= 0),
    product_status product_status NOT NULL DEFAULT 'DRAFT',
    sold_count INTEGER NOT NULL DEFAULT 0 CHECK (sold_count >= 0),
    rating_avg NUMERIC(3, 2) NOT NULL DEFAULT 0 CHECK (rating_avg >= 0 AND rating_avg <= 5),
    total_reviews INTEGER NOT NULL DEFAULT 0 CHECK (total_reviews >= 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    delete_at TIMESTAMPTZ DEFAULT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Foreign key
    CONSTRAINT fk_products_category_id FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);

-- 3. Tạo function cập nhật updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 4. Gắn trigger vào bảng products
CREATE TRIGGER set_updated_at
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_updated_at ON products;

DROP FUNCTION IF EXISTS update_updated_at_column;

DROP TABLE IF EXISTS products;

DROP TYPE IF EXISTS product_status;
-- +goose StatementEnd
