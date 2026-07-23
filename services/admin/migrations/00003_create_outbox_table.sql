-- +goose Up
CREATE TABLE outbox (
    id TEXT PRIMARY KEY,
    topic TEXT NOT NULL,
    key TEXT NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    published BOOLEAN DEFAULT FALSE,
    locked_until TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    retry_count INT DEFAULT 0
);

ALTER TABLE outbox
ADD COLUMN published_at TIMESTAMP,
ADD COLUMN failed BOOLEAN DEFAULT FALSE,
ADD COLUMN last_error TEXT;

-- +goose Down
DROP TABLE IF EXISTS outbox;