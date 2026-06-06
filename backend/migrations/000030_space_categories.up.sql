CREATE TABLE space_categories (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL UNIQUE,
    color       VARCHAR(20),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

ALTER TABLE spaces ADD COLUMN IF NOT EXISTS category_id UUID REFERENCES space_categories(id) ON DELETE SET NULL;

CREATE INDEX idx_spaces_category_id ON spaces (category_id);
