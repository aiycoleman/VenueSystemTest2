-- Filename: migrations/000004_create_venues_table.up.sql
CREATE TABLE IF NOT EXISTS venues (
    id bigserial PRIMARY KEY,
    owner int NOT NULL,
    description text NOT NULL,
    location text NOT NULL,
    price_per_hour DECIMAL(10,2) NOT NULL,
    max_capacity INT NOT NULL,
    image_link text NOT NULL,
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (owner) REFERENCES users(id) ON DELETE CASCADE
);
