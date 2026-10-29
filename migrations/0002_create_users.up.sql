CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    username VARCHAR(32) NOT NULL,
    email VARCHAR(254) NOT NULL,
    password_hash TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT users_username_length_check
        CHECK (char_length(username) BETWEEN 3 AND 32),

    CONSTRAINT users_username_format_check
        CHECK (username ~ '^[a-zA-Z0-9_]+$'),

    CONSTRAINT users_email_length_check
        CHECK (char_length(email) BETWEEN 1 AND 254),

    CONSTRAINT users_username_unique UNIQUE (username),
    CONSTRAINT users_email_unique UNIQUE (email)
);
