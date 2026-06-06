-- =============================================================================
-- PROJECT COMPONENTS — frontend, backend, mobile va h.k.
-- =============================================================================
CREATE TABLE project_components (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id   UUID         NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name         VARCHAR(100) NOT NULL,
    description  TEXT,
    lead_id      UUID         REFERENCES users(id) ON DELETE SET NULL,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ,

    UNIQUE (project_id, name)
);

CREATE INDEX idx_components_project_id ON project_components (project_id) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_components_updated_at
    BEFORE UPDATE ON project_components
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Issue ↔ Component many-to-many
CREATE TABLE issue_components (
    issue_id      UUID NOT NULL REFERENCES issues(id)            ON DELETE CASCADE,
    component_id  UUID NOT NULL REFERENCES project_components(id) ON DELETE CASCADE,
    PRIMARY KEY (issue_id, component_id)
);

CREATE INDEX idx_issue_components_component ON issue_components (component_id);
