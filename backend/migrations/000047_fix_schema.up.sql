-- =============================================================================
-- 000047: Fix schema mismatches between repository code and DB
-- =============================================================================

-- 1. comment_reactions — migration 000044 file exists but table was never
--    applied to this DB instance.
CREATE TABLE IF NOT EXISTS comment_reactions (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    comment_id  UUID        NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
    user_id     UUID        NOT NULL REFERENCES users(id)    ON DELETE CASCADE,
    emoji       VARCHAR(10) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(comment_id, user_id, emoji)
);
CREATE INDEX IF NOT EXISTS idx_comment_reactions_comment ON comment_reactions(comment_id);

-- 2. page_locks — repository code inserts/selects `created_at` but the original
--    migration (000019) created the column as `locked_at`.
--    Add created_at and populate from locked_at so existing rows are consistent.
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'page_locks' AND column_name = 'created_at'
    ) THEN
        ALTER TABLE page_locks ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
        UPDATE page_locks SET created_at = locked_at;
    END IF;
END $$;

-- 3. page_macros — repository code uses column name `macro_type` but the
--    original migration (000019) created the column as `type`.
--    Rename the column so the code and schema agree.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'page_macros' AND column_name = 'type'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'page_macros' AND column_name = 'macro_type'
    ) THEN
        ALTER TABLE page_macros RENAME COLUMN type TO macro_type;
    END IF;
END $$;

-- 4. api_keys — key_prefix is VARCHAR(10) but the usecase generates a 12-char
--    prefix ("jfk_" + 8 hex chars). Widen to 16 to give comfortable headroom.
ALTER TABLE api_keys ALTER COLUMN key_prefix TYPE VARCHAR(16);
