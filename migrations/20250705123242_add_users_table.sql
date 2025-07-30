-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    hashed_pass VARCHAR(255) NOT NULL,
    user_role VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    refresh_token VARCHAR(300) UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
