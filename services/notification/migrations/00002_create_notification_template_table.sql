-- +goose Up
CREATE TABLE IF NOT EXISTS notification_templates (
    id TEXT PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    subject VARCHAR(255),
    body TEXT NOT NULL,
    channel VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS notification_templates;