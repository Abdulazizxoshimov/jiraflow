-- Full-text search upgrade: add tsvector column to issues and pages
-- Note: if search_vector is already a GENERATED column skip ADD/UPDATE.

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'issues' AND column_name = 'search_vector'
  ) THEN
    ALTER TABLE issues ADD COLUMN search_vector tsvector;
  END IF;
END $$;

-- Only update if column is NOT generated (attgenerated = '' means not generated)
DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM pg_attribute a
    JOIN pg_class c ON c.oid = a.attrelid
    WHERE c.relname = 'issues' AND a.attname = 'search_vector' AND a.attgenerated = ''
  ) THEN
    UPDATE issues SET search_vector = to_tsvector('english', coalesce(title, '') || ' ' || coalesce(description, ''));
  END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_issues_search_vector ON issues USING gin(search_vector);

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'pages' AND column_name = 'search_vector'
  ) THEN
    ALTER TABLE pages ADD COLUMN search_vector tsvector;
  END IF;
END $$;

DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM pg_attribute a
    JOIN pg_class c ON c.oid = a.attrelid
    WHERE c.relname = 'pages' AND a.attname = 'search_vector' AND a.attgenerated = ''
  ) THEN
    UPDATE pages SET search_vector = to_tsvector('english', coalesce(title, '') || ' ' || coalesce(content_text, ''));
  END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_pages_search_vector ON pages USING gin(search_vector);
