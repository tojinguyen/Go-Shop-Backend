-- +goose Up
-- +goose StatementBegin
CREATE TYPE promotion_type AS ENUM ('PERCENTAGE', 'VALUE');
CREATE TYPE promotion_status AS ENUM ('DRAFT', 'ACTIVE', 'INACTIVE', 'DELETED');

CREATE TABLE shop_promotions (
    id UUID PRIMARY KEY,
    shop_id UUID NOT NULL,
    promotion_name VARCHAR(255) NOT NULL,
    promotion_type promotion_type NOT NULL,
    discount_value DECIMAL(10,2) NOT NULL CHECK (discount_value > 0),
    max_discount_amount DECIMAL(10,2) CHECK (max_discount_amount > 0),
    min_purchase_amount DECIMAL(10,2) DEFAULT 0.00 CHECK (min_purchase_amount >= 0),
    usage_limit_per_user INTEGER DEFAULT 1 CHECK (usage_limit_per_user > 0),
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    promotion_status promotion_status DEFAULT 'DRAFT',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Add constraint to ensure start_time is before end_time
    CONSTRAINT chk_promotion_time_range CHECK (start_time < end_time),
    
    -- Add constraint for percentage type (should be between 0 and 100)
    CONSTRAINT chk_percentage_value CHECK (
        (promotion_type = 'PERCENTAGE' AND discount_value <= 100) OR 
        (promotion_type = 'VALUE')
    )
);


ALTER TABLE shop_promotions ADD CONSTRAINT fk_shop_promotions_shop_id 
    FOREIGN KEY (shop_id) REFERENCES shops(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS shop_promotions;
DROP TYPE IF EXISTS promotion_status;
DROP TYPE IF EXISTS promotion_type;
-- +goose StatementEnd
