CREATE TABLE IF NOT EXISTS refresh_tokens (
    id bytea PRIMARY KEY,
    token_value VARCHAR(255) NOT NULL,
    user_id bytea NOT NULL,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    revoked BOOLEAN DEFAULT FALSE
);
