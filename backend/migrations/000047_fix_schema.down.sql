-- Reverse of 000047_fix_schema.up.sql

DROP TABLE IF EXISTS comment_reactions;

ALTER TABLE page_locks DROP COLUMN IF EXISTS created_at;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'page_macros' AND column_name = 'macro_type'
    ) THEN
        ALTER TABLE page_macros RENAME COLUMN macro_type TO type;
    END IF;
END $$;

ALTER TABLE api_keys ALTER COLUMN key_prefix TYPE VARCHAR(10);
