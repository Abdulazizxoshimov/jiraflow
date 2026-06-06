-- =============================================================================
-- PROJECTS — loyihalar
-- =============================================================================
CREATE TABLE projects (
    id             UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    key            VARCHAR(10)  NOT NULL UNIQUE
                   CHECK (key ~ '^[A-Z][A-Z0-9]{1,9}$'),  -- PROJ, ABC, WEB
    name           VARCHAR(255) NOT NULL,
    description    TEXT,
    icon_url       TEXT,
    lead_id        UUID         NOT NULL REFERENCES users(id),
    workflow_id    UUID         NOT NULL REFERENCES workflows(id),
    issue_counter  BIGINT       NOT NULL DEFAULT 0,  -- keyingi PROJ-N raqami
    is_archived    BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at     TIMESTAMPTZ
);

CREATE INDEX idx_projects_lead_id     ON projects (lead_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_projects_is_archived ON projects (is_archived) WHERE deleted_at IS NULL;
CREATE INDEX idx_projects_name_trgm   ON projects USING GIN (name gin_trgm_ops);

CREATE TRIGGER trg_projects_updated_at
    BEFORE UPDATE ON projects
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMENT ON COLUMN projects.key IS 'Unique kalit: 2-10 katta harf/raqam (PROJ, ABC123)';
COMMENT ON COLUMN projects.issue_counter IS 'Atomically increment qilib PROJ-N raqamini olamiz';

-- =============================================================================
-- PROJECT_MEMBERS — loyihaga foydalanuvchilar va ularning rollari
-- =============================================================================
CREATE TABLE project_members (
    project_id   UUID         NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id      UUID         NOT NULL REFERENCES users(id)    ON DELETE CASCADE,
    role         VARCHAR(32)  NOT NULL DEFAULT 'member'
                 CHECK (role IN ('admin', 'member', 'viewer')),
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    PRIMARY KEY (project_id, user_id)
);

CREATE INDEX idx_project_members_user_id ON project_members (user_id);

-- =============================================================================
-- CUSTOM_FIELDS — loyiha darajasida custom field'lar
-- =============================================================================
CREATE TABLE custom_fields (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id   UUID         NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name         VARCHAR(100) NOT NULL,
    field_key    VARCHAR(50)  NOT NULL,  -- API'da ishlatish uchun, slug
    field_type   VARCHAR(32)  NOT NULL
                 CHECK (field_type IN ('text', 'number', 'date', 'select', 'multi_select', 'user', 'checkbox', 'url')),
    is_required  BOOLEAN      NOT NULL DEFAULT FALSE,
    options      JSONB,       -- select uchun: ["Bug", "Feature", "Task"]
    position     INT          NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ,

    UNIQUE (project_id, field_key)
);

CREATE INDEX idx_custom_fields_project_id ON custom_fields (project_id) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_custom_fields_updated_at
    BEFORE UPDATE ON custom_fields
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
