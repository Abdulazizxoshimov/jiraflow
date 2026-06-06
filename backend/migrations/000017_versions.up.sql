-- =============================================================================
-- PROJECT VERSIONS / RELEASES — v1.0, v2.3-beta va h.k.
-- =============================================================================
CREATE TABLE project_versions (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id   UUID         NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name         VARCHAR(100) NOT NULL,
    description  TEXT,
    status       VARCHAR(16)  NOT NULL DEFAULT 'unreleased'
                 CHECK (status IN ('unreleased', 'released', 'archived')),
    start_date   DATE,
    release_date DATE,
    released_at  TIMESTAMPTZ,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ,

    UNIQUE (project_id, name)
);

CREATE INDEX idx_versions_project_id ON project_versions (project_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_versions_status     ON project_versions (project_id, status) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_versions_updated_at
    BEFORE UPDATE ON project_versions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Issue ↔ Version many-to-many (fix version)
CREATE TABLE issue_versions (
    issue_id    UUID NOT NULL REFERENCES issues(id)          ON DELETE CASCADE,
    version_id  UUID NOT NULL REFERENCES project_versions(id) ON DELETE CASCADE,
    PRIMARY KEY (issue_id, version_id)
);

CREATE INDEX idx_issue_versions_version ON issue_versions (version_id);
