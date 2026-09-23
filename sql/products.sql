CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    price NUMERIC(10, 2) NOT NULL
);

INSERT INTO products (name, price) VALUES ('Kopi Hitam', 15000.00);

ALTER TABLE products ADD stocks INT DEFAULT 0;
ALTER TABLE products ADD status VARCHAR(20) DEFAULT 'ACTIVE';