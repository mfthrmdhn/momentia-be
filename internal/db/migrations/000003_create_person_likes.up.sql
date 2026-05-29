CREATE TABLE person_likes (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    person_id        UUID NOT NULL REFERENCES person(id) ON DELETE CASCADE,
    liked_food       TEXT,
    liked_places     TEXT,
    liked_color      VARCHAR(100),
    additional_notes TEXT,
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    updated_at       TIMESTAMPTZ DEFAULT NOW()
);
