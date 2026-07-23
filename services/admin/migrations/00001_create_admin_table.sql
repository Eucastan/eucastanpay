-- +goose Up
CREATE TYPE admin_role AS ENUM ('super_admin', 'admin', 'moderator');

CREATE TABLE IF NOT EXISTS admins (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,

    first_name TEXT,
    last_name TEXT,

    role admin_role DEFAULT 'admin',
    status TEXT DEFAULT 'active',

    two_fa_enabled BOOLEAN DEFAULT FALSE,
    two_fa_secret TEXT,

    last_login_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_admins_email ON admins(email);

-- +goose Down
DROP TABLE IF EXISTS admins;
DROP TYPE IF EXISTS admin_role;