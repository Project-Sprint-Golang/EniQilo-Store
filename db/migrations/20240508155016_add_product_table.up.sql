CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    sku VARCHAR(255),
    category VARCHAR(255) CHECK (category IN ('Clothing', 'Accessories', 'Footwear', 'Beverages')),
    imageUrl TEXT[],
    notes TEXT,
    price DECIMAL(10, 2) NOT NULL,
    stock INT NOT NULL,
    location VARCHAR(255),
    isAvailable BOOLEAN NOT NULL DEFAULT true,
    createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deletedAt TIMESTAMP
);

CREATE INDEX idx_products_id ON products (id)