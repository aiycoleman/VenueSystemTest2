-- Filename: migrations/000001_create_roles_table.up.sql
CREATE TABLE IF NOT EXISTS roles (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);