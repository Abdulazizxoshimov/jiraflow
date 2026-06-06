-- =============================================================================
-- FAZA 4 — Ilg'or xususiyatlar
-- =============================================================================

-- PAGE REACTIONS (emoji) ─────────────────────────────────────────────────────
CREATE TABLE page_reactions (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    page_id    UUID        NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    emoji      VARCHAR(32) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (page_id, user_id, emoji)
);

CREATE INDEX idx_page_reactions_page ON page_reactions (page_id, emoji);

-- ISSUE MULTIPLE ASSIGNEES ───────────────────────────────────────────────────
-- issue_assignees — primary assignee + co-assignees
CREATE TABLE issue_assignees (
    issue_id   UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (issue_id, user_id)
);

CREATE INDEX idx_issue_assignees_user ON issue_assignees (user_id);

-- WEBHOOKS ───────────────────────────────────────────────────────────────────
CREATE TABLE webhooks (
    id          UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  UUID          REFERENCES projects(id) ON DELETE CASCADE,
    space_id    UUID          REFERENCES spaces(id) ON DELETE CASCADE,
    name        VARCHAR(255)  NOT NULL,
    url         TEXT          NOT NULL,
    secret      VARCHAR(255),
    events      TEXT[]        NOT NULL DEFAULT '{}',  -- issue.created, issue.updated, ...
    is_active   BOOLEAN       NOT NULL DEFAULT TRUE,
    created_by  UUID          NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),

    CHECK (project_id IS NOT NULL OR space_id IS NOT NULL)
);

CREATE INDEX idx_webhooks_project ON webhooks (project_id) WHERE project_id IS NOT NULL AND is_active;
CREATE INDEX idx_webhooks_space   ON webhooks (space_id) WHERE space_id IS NOT NULL AND is_active;

CREATE TRIGGER trg_webhooks_updated_at
    BEFORE UPDATE ON webhooks
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- WEBHOOK DELIVERIES — delivery log ─────────────────────────────────────────
CREATE TABLE webhook_deliveries (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    webhook_id    UUID        NOT NULL REFERENCES webhooks(id) ON DELETE CASCADE,
    event         VARCHAR(64) NOT NULL,
    payload       JSONB       NOT NULL,
    status_code   INT,
    response_body TEXT,
    success       BOOLEAN     NOT NULL DEFAULT FALSE,
    delivered_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_webhook_deliveries_webhook ON webhook_deliveries (webhook_id, delivered_at DESC);

-- PAGE EDITING LOCK (collaborative editing minimal versiya) ──────────────────
CREATE TABLE page_locks (
    page_id    UUID        PRIMARY KEY REFERENCES pages(id) ON DELETE CASCADE,
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_id VARCHAR(64) NOT NULL,
    locked_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_page_locks_expires ON page_locks (expires_at);

-- PAGE MACROS (Jira issue embed) ─────────────────────────────────────────────
CREATE TABLE page_macros (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    page_id     UUID        NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    type        VARCHAR(32) NOT NULL CHECK (type IN ('jira_issues', 'jira_sprint', 'status_page')),
    config      JSONB       NOT NULL DEFAULT '{}'::jsonb,
    position    INT         NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_page_macros_page ON page_macros (page_id, position);

-- LEXORANK: issues.rank column ───────────────────────────────────────────────
ALTER TABLE issues ADD COLUMN IF NOT EXISTS rank VARCHAR(255) NOT NULL DEFAULT '';
CREATE INDEX idx_issues_rank ON issues (project_id, rank) WHERE deleted_at IS NULL;
