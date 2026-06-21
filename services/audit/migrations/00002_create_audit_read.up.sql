-- +goose Up
CREATE TABLE IF NOT EXISTS audit_read (
    id TEXT PRIMARY KEY,
    event_type TEXT NOT NULL,
    service TEXT NOT NULL,
    correlation_id TEXT,
    causation_id TEXT,
    reference TEXT,
    account_id TEXT,
    user_id TEXT,
    amount BIGINT,
    status TEXT,
    payload JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_read_corr ON audit_read(correlation_id);
CREATE INDEX IF NOT EXISTS idx_audit_read_ref ON audit_read(reference);
CREATE INDEX IF NOT EXISTS idx_audit_read_event ON audit_read(event_type);
CREATE INDEX IF NOT EXISTS idx_audit_read_created ON audit_read(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_read_account ON audit_read(account_id);

-- +goose Down
DROP TABLE IF EXISTS audit_read;
DROP INDEX IF EXISTS idx_audit_read_corr;
DROP INDEX IF EXISTS idx_audit_read_ref;
DROP INDEX IF EXISTS idx_audit_read_event;
DROP INDEX IF EXISTS idx_audit_read_created;
DROP INDEX IF EXISTS idx_audit_read_account;