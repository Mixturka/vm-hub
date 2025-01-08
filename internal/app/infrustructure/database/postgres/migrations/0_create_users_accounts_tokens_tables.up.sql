-- 1__create_users_accounts_tokens_tables.up.sql

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY,
    profile_picture TEXT DEFAULT '',
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT UNIQUE NOT NULL,
    is_email_verified BOOLEAN DEFAULT FALSE,
    is_two_factor_enabled BOOLEAN DEFAULT FALSE,
    method INTEGER NOT NULL, -- AuthMethod enum (0: Credentials, 1: Google, 2: Yandex)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    type TEXT NOT NULL,
    provider TEXT NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE, -- Foreign key to users
    refresh_token TEXT,
    access_token TEXT,
    expires_at INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tokens (
    id UUID PRIMARY KEY,
    user_email TEXT NOT NULL,
    token TEXT UNIQUE NOT NULL,
    type INTEGER NOT NULL, -- TokenType enum (0: Verification, 1: TwoFactor, 2: PasswordReset)
    epires_in TIMESTAMP NOT NULL
);
