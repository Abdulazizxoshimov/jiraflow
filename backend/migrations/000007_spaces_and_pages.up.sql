-- =============================================================================
-- SPACES — Knowledge Base spaces
-- =============================================================================
CREATE TABLE spaces (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    key           VARCHAR(10)  NOT NULL UNIQUE
                  CHECK (key ~ '^[A-Z][A-Z0-9]{1,9}$'),
    name          VARCHAR(255) NOT NULL,
    description   TEXT,
    icon_url      TEXT,
    type          VARCHAR(32)  NOT NULL DEFAULT 'team'
                  CHECK (type IN ('team', 'personal', 'project')),
    lead_id       UUID         NOT NULL REFERENCES users(id),
    project_id    UUID         REFERENCES projects(id) ON DELETE SET NULL,  -- type='project' uchun
    is_archived   BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

CREATE INDEX idx_spaces_lead_id     ON spaces (lead_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_spaces_project_id  ON spaces (project_id) WHERE deleted_at IS NULL AND project_id IS NOT NULL;
CREATE INDEX idx_spaces_is_archived ON spaces (is_archived) WHERE deleted_at IS NULL;
CREATE INDEX idx_spaces_name_trgm   ON spaces USING GIN (name gin_trgm_ops);

CREATE TRIGGER trg_spaces_updated_at
    BEFORE UPDATE ON spaces
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- =============================================================================
-- SPACE_MEMBERS — space'ga foydalanuvchilar
-- =============================================================================
CREATE TABLE space_members (
    space_id    UUID         NOT NULL REFERENCES spaces(id) ON DELETE CASCADE,
    user_id     UUID         NOT NULL REFERENCES users(id)  ON DELETE CASCADE,
    role        VARCHAR(32)  NOT NULL DEFAULT 'member'
                CHECK (role IN ('admin', 'member', 'viewer')),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    PRIMARY KEY (space_id, user_id)
);

CREATE INDEX idx_space_members_user_id ON space_members (user_id);

-- =============================================================================
-- PAGES — Confluence sahifalari (ierarxik)
-- =============================================================================
CREATE TABLE pages (
    id              UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    space_id        UUID         NOT NULL REFERENCES spaces(id) ON DELETE CASCADE,
    parent_id       UUID         REFERENCES pages(id) ON DELETE CASCADE,
    title           VARCHAR(500) NOT NULL,
    -- Rich content (ProseMirror/TipTap JSON yoki sanitized HTML)
    content         JSONB        NOT NULL DEFAULT '{}'::jsonb,
    -- Search uchun ajratilgan plain text
    content_text    TEXT         NOT NULL DEFAULT '',

    author_id       UUID         NOT NULL REFERENCES users(id),
    last_editor_id  UUID         NOT NULL REFERENCES users(id),
    current_version INT          NOT NULL DEFAULT 1,
    status          VARCHAR(16)  NOT NULL DEFAULT 'published'
                    CHECK (status IN ('draft', 'published')),
    position        INT          NOT NULL DEFAULT 0,  -- tree'da tartib

    search_vector   tsvector     GENERATED ALWAYS AS (
        setweight(to_tsvector('simple', coalesce(title, '')), 'A') ||
        setweight(to_tsvector('simple', coalesce(content_text, '')), 'B')
    ) STORED,

    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,

    CHECK (id <> parent_id)
);

CREATE INDEX idx_pages_space_parent ON pages (space_id, parent_id, position) WHERE deleted_at IS NULL;
CREATE INDEX idx_pages_author       ON pages (author_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_pages_status       ON pages (space_id, status) WHERE deleted_at IS NULL;
CREATE INDEX idx_pages_search       ON pages USING GIN (search_vector);
CREATE INDEX idx_pages_title_trgm   ON pages USING GIN (title gin_trgm_ops);

CREATE TRIGGER trg_pages_updated_at
    BEFORE UPDATE ON pages
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMENT ON COLUMN pages.content IS 'TipTap/ProseMirror JSON schema';
COMMENT ON COLUMN pages.content_text IS 'Frontend ajratib jo''natadi yoki backend extract qiladi';

-- =============================================================================
-- PAGE_VERSIONS — sahifa versiyalari (history)
-- =============================================================================
CREATE TABLE page_versions (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    page_id       UUID         NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    version       INT          NOT NULL,
    title         VARCHAR(500) NOT NULL,
    content       JSONB        NOT NULL,
    content_text  TEXT         NOT NULL DEFAULT '',
    author_id     UUID         NOT NULL REFERENCES users(id),
    change_note   TEXT,                       -- ixtiyoriy commit-message
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    UNIQUE (page_id, version)
);

CREATE INDEX idx_page_versions_page_id ON page_versions (page_id, version DESC);

-- =============================================================================
-- PAGE_WATCHERS
-- =============================================================================
CREATE TABLE page_watchers (
    page_id     UUID         NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    user_id     UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    PRIMARY KEY (page_id, user_id)
);

CREATE INDEX idx_page_watchers_user_id ON page_watchers (user_id);
