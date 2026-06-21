-- +goose Up
CREATE TABLE IF NOT EXISTS accounts (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    account_no BIGINT NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0 CHECK (balance >= 0),
    account_type TEXT NOT NULL,
    currency TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_accounts_user_account_type ON accounts (user_id, account_type);
CREATE UNIQUE INDEX idx_accounts_account_no ON accounts (account_no);

-- +goose Down
DROP INDEX IF EXISTS idx_accounts_account_no;
DROP INDEX IF EXISTS idx_accounts_user_account_type;
DROP TABLE IF EXISTS accounts;