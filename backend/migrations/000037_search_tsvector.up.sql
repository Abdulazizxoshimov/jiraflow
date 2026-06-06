-- Full-text search upgrade: add tsvector column to issues and pages

ALTER TABLE issues ADD COLUMN IF NOT EXISTS search_vector tsvector;
UPDATE issues SET search_vector = to_tsvector('english', coalesce(title, '') || ' ' || coalesce(description, ''));
CREATE INDEX IF NOT EXISTS idx_issues_search_vector ON issues USING gin(search_vector);

CREATE OR REPLACE FUNCTION issues_search_vector_update() RETURNS trigger AS $$
BEGIN
  NEW.search_vector := to_tsvector('english', coalesce(NEW.title, '') || ' ' || coalesce(NEW.description, ''));
  RETURN NEW;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS issues_search_vector_trig ON issues;
CREATE TRIGGER issues_search_vector_trig
  BEFORE INSERT OR UPDATE OF title, description ON issues
  FOR EACH ROW EXECUTE FUNCTION issues_search_vector_update();

ALTER TABLE pages ADD COLUMN IF NOT EXISTS search_vector tsvector;
UPDATE pages SET search_vector = to_tsvector('english', coalesce(title, '') || ' ' || coalesce(content_text, ''));
CREATE INDEX IF NOT EXISTS idx_pages_search_vector ON pages USING gin(search_vector);

CREATE OR REPLACE FUNCTION pages_search_vector_update() RETURNS trigger AS $$
BEGIN
  NEW.search_vector := to_tsvector('english', coalesce(NEW.title, '') || ' ' || coalesce(NEW.content_text, ''));
  RETURN NEW;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS pages_search_vector_trig ON pages;
CREATE TRIGGER pages_search_vector_trig
  BEFORE INSERT OR UPDATE OF title, content_text ON pages
  FOR EACH ROW EXECUTE FUNCTION pages_search_vector_update();
