-- Filename: migrations/000003_create_users_table.up.sql
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    email citext UNIQUE NOT NULL,
    role int NOT NULL,
    password_hash bytea NOT NULL,
    activated bool NOT NULL DEFAULT TRUE,
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (role) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT user_email_key UNIQUE (email)
);