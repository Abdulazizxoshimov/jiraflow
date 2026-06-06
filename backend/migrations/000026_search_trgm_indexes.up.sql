-- pg_trgm extension 000001 da qo'shilgan, bu yerda faqat GIN indexlar
CREATE INDEX IF NOT EXISTS idx_issues_title_trgm
    ON issues USING gin(title gin_trgm_ops)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_pages_title_trgm
    ON pages USING gin(title gin_trgm_ops)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_projects_name_trgm
    ON projects USING gin(name gin_trgm_ops)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_spaces_name_trgm
    ON spaces USING gin(name gin_trgm_ops)
    WHERE deleted_at IS NULL;
