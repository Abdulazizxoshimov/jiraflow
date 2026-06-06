DROP TRIGGER IF EXISTS issues_search_vector_trig ON issues;
DROP FUNCTION IF EXISTS issues_search_vector_update();
DROP INDEX IF EXISTS idx_issues_search_vector;
ALTER TABLE issues DROP COLUMN IF EXISTS search_vector;

DROP TRIGGER IF EXISTS pages_search_vector_trig ON pages;
DROP FUNCTION IF EXISTS pages_search_vector_update();
DROP INDEX IF EXISTS idx_pages_search_vector;
ALTER TABLE pages DROP COLUMN IF EXISTS search_vector;
