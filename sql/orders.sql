CREATE TABLE orders (
    id VARCHAR(50) PRIMARY KEY,
    user_id int NOT NULL,
    product_id int NOT NULL,
    quantity int NOT NULL,
    total_amount NUMERIC(15, 2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
)