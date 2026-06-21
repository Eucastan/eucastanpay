
-- +goose Up
CREATE TABLE IF NOT EXISTS processed_events (
    id TEXT PRIMARY KEY,
    event_id TEXT UNIQUE,
    topic TEXT,
    processed_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS processed_events;