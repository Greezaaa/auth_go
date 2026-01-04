-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    name TEXT NOT NULL,
    allowed_apps TEXT[],
    is_admin BOOLEAN DEFAULT false,
    is_confirmed BOOLEAN DEFAULT false,
    is_email_confirmed BOOLEAN DEFAULT false,
    email_code TEXT,
    role TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE users;
DROP EXTENSION IF EXISTS "uuid-ossp";