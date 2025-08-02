-- +goose Up
CREATE TABLE users (
                       id BIGSERIAL PRIMARY KEY,
                       telegram_id BIGINT UNIQUE NOT NULL,
                       username VARCHAR(64),
                       password VARCHAR(255),
                       created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE users;