-- +goose Up
CREATE TABLE IF NOT EXISTS transfers(
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    reference TEXT NOT NULL UNIQUE,
    step VARCHAR(50),
    from_account_id TEXT NOT NULL,
    from_account_no BIGINT NOT NULL,
    to_account_id TEXT NOT NULL,
    to_account_no BIGINT NOT NULL,
    amount BIGINT NOT NULL CHECK (amount > 0),
    description TEXT,
    idempotency_key TEXT UNIQUE NOT NULL,
    type TEXT NOT NULL,         
	status TEXT NOT NULL,
	mode TEXT NOT NULL,
	reversal_ref TEXT NOT NULL DEFAULT '',
	is_reversed BOOLEAN DEFAULT FALSE,
    from_balance_after BIGINT,
    to_balance_after BIGINT,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE transfers
ADD COLUMN recovery_count INT DEFAULT 0,
ADD COLUMN last_recovery_at TIMESTAMP;

CREATE INDEX idx_transfer_reference ON transfers(reference);
CREATE INDEX idx_transfer_from_account ON transfers(from_account_id);
CREATE INDEX idx_transfer_to_account ON transfers(to_account_id);
CREATE INDEX idx_transfer_user ON transfers(user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_transfers_idempotency_key ON transfers(idempotency_key);
CREATE INDEX IF NOT EXISTS idx_transfers_status_step ON transfers(status, step);
CREATE INDEX IF NOT EXISTS idx_transfers_created_at ON transfers(created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS transfers;
DROP INDEX IF EXISTS idx_transfers_idempotency_key;
DROP INDEX IF EXISTS idx_transfers_status_step;
DROP INDEX IF EXISTS idx_transfers_created_at;