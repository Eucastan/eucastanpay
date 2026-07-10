
-- +goose Up
CREATE TABLE IF NOT EXISTS ledgers (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    account_id TEXT NOT NULL,
    amount BIGINT NOT NULL CHECK (amount > 0),
    entry_type TEXT NOT NULL,  
    reference TEXT NOT NULL,
    balance_after BIGINT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ledgers_account_id ON ledgers(account_id);
CREATE INDEX IF NOT EXISTS idx_ledgers_reference ON ledgers(reference);
CREATE INDEX IF NOT EXISTS idx_ledgers_created_at ON ledgers(created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS ledgers;
DROP INDEX IF EXISTS idx_ledgers_account_id;
DROP INDEX IF EXISTS idx_ledgers_reference;
DROP INDEX IF EXISTS idx_ledgers_created_at;