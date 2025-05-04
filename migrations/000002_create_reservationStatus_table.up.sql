-- Filename: migrations/000002_create_reservationStatus_table.up.sql
CREATE TABLE IF NOT EXISTS reservationStatus (
    id bigserial PRIMARY KEY,
    status text NOT NULL,
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);