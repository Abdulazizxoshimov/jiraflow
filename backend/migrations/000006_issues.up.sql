-- =============================================================================
-- ISSUES — vazifalar (asosiy jadval)
-- =============================================================================
CREATE TABLE issues (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id    UUID         NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    issue_number  INT          NOT NULL,  -- PROJ-N dagi N
    title         VARCHAR(500) NOT NULL,
    description   TEXT,        -- markdown yoki JSON (TipTap)

    type          VARCHAR(32)  NOT NULL DEFAULT 'task'
                  CHECK (type IN ('task', 'bug', 'story', 'epic', 'subtask')),
    status_id     UUID         NOT NULL REFERENCES workflow_statuses(id),
    priority      VARCHAR(16)  NOT NULL DEFAULT 'medium'
                  CHECK (priority IN ('lowest', 'low', 'medium', 'high', 'highest')),

    assignee_id   UUID         REFERENCES users(id) ON DELETE SET NULL,
    reporter_id   UUID         NOT NULL REFERENCES users(id),
    parent_id     UUID         REFERENCES issues(id) ON DELETE SET NULL,  -- subtask uchun
    sprint_id     UUID         REFERENCES sprints(id) ON DELETE SET NULL,

    story_points  INT          CHECK (story_points IS NULL OR story_points >= 0),
    due_date      DATE,

    -- Custom field qiymatlari: { "field_key": value, ... }
    custom_fields JSONB        NOT NULL DEFAULT '{}'::jsonb,

    -- Full-text search uchun avtomatik yangilanadigan ustun
    search_vector tsvector     GENERATED ALWAYS AS (
        setweight(to_tsvector('simple', coalesce(title, '')), 'A') ||
        setweight(to_tsvector('simple', coalesce(description, '')), 'B')
    ) STORED,

    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ,

    UNIQUE (project_id, issue_number)
);

-- Asosiy indekslar
CREATE INDEX idx_issues_project_status ON issues (project_id, status_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_issues_assignee       ON issues (assignee_id) WHERE deleted_at IS NULL AND assignee_id IS NOT NULL;
CREATE INDEX idx_issues_reporter       ON issues (reporter_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_issues_sprint         ON issues (sprint_id) WHERE deleted_at IS NULL AND sprint_id IS NOT NULL;
CREATE INDEX idx_issues_parent         ON issues (parent_id) WHERE deleted_at IS NULL AND parent_id IS NOT NULL;
CREATE INDEX idx_issues_priority       ON issues (project_id, priority) WHERE deleted_at IS NULL;
CREATE INDEX idx_issues_created_at     ON issues (created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_issues_due_date       ON issues (due_date) WHERE deleted_at IS NULL AND due_date IS NOT NULL;
CREATE INDEX idx_issues_custom_fields  ON issues USING GIN (custom_fields);
CREATE INDEX idx_issues_search         ON issues USING GIN (search_vector);

CREATE TRIGGER trg_issues_updated_at
    BEFORE UPDATE ON issues
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMENT ON COLUMN issues.search_vector IS 'Avtomatik to''ldiriladi title+description bo''yicha';

-- =============================================================================
-- LABELS — teglar
-- =============================================================================
CREATE TABLE labels (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id   UUID         NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name         VARCHAR(64)  NOT NULL,
    color        VARCHAR(7)   NOT NULL DEFAULT '#6B7280',
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    UNIQUE (project_id, name)
);

CREATE INDEX idx_labels_project_id ON labels (project_id);

-- =============================================================================
-- ISSUE_LABELS — many-to-many vazifa va teg
-- =============================================================================
CREATE TABLE issue_labels (
    issue_id  UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    label_id  UUID NOT NULL REFERENCES labels(id) ON DELETE CASCADE,
    PRIMARY KEY (issue_id, label_id)
);

CREATE INDEX idx_issue_labels_label_id ON issue_labels (label_id);

-- =============================================================================
-- ISSUE_WATCHERS — vazifani kuzatuvchilar
-- =============================================================================
CREATE TABLE issue_watchers (
    issue_id    UUID         NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    user_id     UUID         NOT NULL REFERENCES users(id)  ON DELETE CASCADE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    PRIMARY KEY (issue_id, user_id)
);

CREATE INDEX idx_issue_watchers_user_id ON issue_watchers (user_id);

-- =============================================================================
-- ISSUE_LINKS — vazifalar orasidagi bog'lanishlar
-- =============================================================================
CREATE TABLE issue_links (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id     UUID         NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    target_id     UUID         NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    link_type     VARCHAR(32)  NOT NULL
                  CHECK (link_type IN ('relates_to', 'blocks', 'blocked_by', 'duplicates', 'is_duplicated_by')),
    created_by    UUID         NOT NULL REFERENCES users(id),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CHECK (source_id <> target_id),
    UNIQUE (source_id, target_id, link_type)
);

CREATE INDEX idx_issue_links_source ON issue_links (source_id);
CREATE INDEX idx_issue_links_target ON issue_links (target_id);

-- =============================================================================
-- ISSUE_HISTORY — vazifa o'zgarishlar tarixi (audit log)
-- =============================================================================
CREATE TABLE issue_history (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    issue_id    UUID         NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    user_id     UUID         REFERENCES users(id) ON DELETE SET NULL,
    field       VARCHAR(64)  NOT NULL,  -- 'status', 'assignee', 'title', va h.k.
    old_value   JSONB,                  -- eski qiymat
    new_value   JSONB,                  -- yangi qiymat
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_issue_history_issue_id ON issue_history (issue_id, created_at DESC);
