CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    phoneNumber VARCHAR(20) NOT NULL,
    name VARCHAR(50) NOT NULL,
    password VARCHAR(100) NOT NULL,
    role INTEGER NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk mempercepat pencarian berdasarkan nomor telepon
CREATE INDEX idx_users_phone_number ON users (phoneNumber);

-- Tabel untuk menyimpan token JWT
CREATE TABLE IF NOT EXISTS JWTToken (
    id SERIAL PRIMARY KEY,
    userid INT ,
    token TEXT ,
    FOREIGN KEY (userId) REFERENCES users(id)
);

-- Untuk mempercepat pencarian berdasarkan nomor telepon, indeks tambahan dibuat.
CREATE INDEX idx_jwt_token_user_id ON JWTToken (userId);