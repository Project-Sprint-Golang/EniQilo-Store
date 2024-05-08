CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    phoneNumber VARCHAR(20) NOT NULL,
    name VARCHAR(50) NOT NULL,
    password VARCHAR(100) NOT NULL,
    role INTEGER NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deletedAt TIMESTAMP
);

-- Index untuk mempercepat pencarian berdasarkan nomor telepon
CREATE INDEX idx_users_phone_number ON users (phoneNumber);
CREATE INDEX idx_users_id ON users (id);

