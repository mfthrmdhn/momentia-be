CREATE TABLE person_dates (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    person_id  UUID NOT NULL REFERENCES person(id) ON DELETE CASCADE,
    label      VARCHAR(100) NOT NULL,
    date       DATE NOT NULL,
    note       VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_person_dates_person_id ON person_dates(person_id);
CREATE INDEX idx_person_dates_date ON person_dates(date);
