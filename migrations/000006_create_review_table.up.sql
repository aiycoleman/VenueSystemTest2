-- Filename: migrations/000006_create_review_table.up.sql
CREATE TABLE IF NOT EXISTS review (
	id bigserial PRIMARY KEY,
	customer int NOT NULL,
	venue int NOT NULL,
	comment text NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (venue) REFERENCES venue(id) ON DELETE CASCADE
);