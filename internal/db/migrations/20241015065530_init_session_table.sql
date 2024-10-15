-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS auth.sessions(
    session_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE,
    user_email VARCHAR(255) NOT NULL,
    refresh_token VARCHAR(512) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
);
-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd
