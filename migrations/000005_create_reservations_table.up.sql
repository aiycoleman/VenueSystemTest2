-- Filename: migrations/000005_create_reservations_table.up.sql
CREATE TABLE IF NOT EXISTS reservations (
    id bigserial PRIMARY KEY,
    venue int NOT NULL,
    customer int NOT NULL,
    start_date date NOT NULL,
    start_time time NOT NULL,  
    end_time time NOT NULL, 
    status int NOT NULL,
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (venue) REFERENCES venues(id) ON DELETE CASCADE,
    FOREIGN KEY (customer) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (status) REFERENCES reservationStatus(id) ON DELETE CASCADE
);
