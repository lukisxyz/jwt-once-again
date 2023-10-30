CREATE TABLE IF NOT EXISTS accounts (
    id BYTEA PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    code_verification VARCHAR(255),
    email_verified_at TIMESTAMPTZ
);
