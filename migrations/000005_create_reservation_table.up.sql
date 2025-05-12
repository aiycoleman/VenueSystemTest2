-- Filename: migrations/000005_create_reservation_table.up.sql
CREATE TABLE IF NOT EXISTS reservation (
    id bigserial PRIMARY KEY,
    venue int NOT NULL,
    customer int NOT NULL,
    start_date date NOT NULL,
    start_time time NOT NULL,  
    end_time time NOT NULL, 
    status int NOT NULL DEFAULT 1,
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (venue) REFERENCES venue(id) ON DELETE CASCADE,
    FOREIGN KEY (customer) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (status) REFERENCES reservationStatus(id) ON DELETE CASCADE
);
