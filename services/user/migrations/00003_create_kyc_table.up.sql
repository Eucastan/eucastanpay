-- +goose Up
CREATE TABLE IF NOT EXISTS kycs (
    id TEXT PRIMARY KEY,
    user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
    id_type TEXT,
    id_number TEXT,
    status TEXT DEFAULT 'pending',
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS kycs;