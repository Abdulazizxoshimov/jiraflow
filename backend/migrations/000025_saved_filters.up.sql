CREATE TABLE saved_filters (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    filter_type VARCHAR(32)  NOT NULL DEFAULT 'issue', -- issue | page
    filters     JSONB        NOT NULL DEFAULT '{}'::jsonb,
    is_shared   BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_saved_filters_user   ON saved_filters (user_id);
CREATE INDEX idx_saved_filters_shared ON saved_filters (is_shared) WHERE is_shared = TRUE;
CREATE INDEX idx_saved_filters_data   ON saved_filters USING GIN (filters);

CREATE TRIGGER trg_saved_filters_updated_at
    BEFORE UPDATE ON saved_filters
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
