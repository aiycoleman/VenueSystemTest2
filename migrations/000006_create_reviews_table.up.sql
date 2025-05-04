-- Filename: migrations/000006_create_reviews_table.up.sql
CREATE TABLE IF NOT EXISTS reviews (
	id bigserial PRIMARY KEY,
	customer int NOT NULL,
	venue int NOT NULL,
	comment text NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (venue) REFERENCES venues(id) ON DELETE CASCADE
);