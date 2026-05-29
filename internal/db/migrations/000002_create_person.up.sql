CREATE TABLE person (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    creator_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,
    relationship    VARCHAR(255) NOT NULL,
    is_pinned       BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);
