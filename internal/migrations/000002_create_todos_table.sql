-- +goose Up
CREATE TABLE todos (
                       id SERIAL PRIMARY KEY,
                       user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                       title VARCHAR(255) NOT NULL,
                       description TEXT,
                       due_date TIMESTAMP,
                       completed BOOLEAN DEFAULT FALSE,
                       created_at TIMESTAMP DEFAULT NOW(),
                       updated_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE todos;