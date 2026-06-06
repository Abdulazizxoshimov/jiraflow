-- =============================================================================
-- BOARDS — Kanban/Scrum board konfiguratsiyasi
-- =============================================================================
CREATE TABLE boards (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id   UUID         NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name         VARCHAR(255) NOT NULL,
    type         VARCHAR(16)  NOT NULL DEFAULT 'kanban'
                 CHECK (type IN ('kanban', 'scrum')),
    -- Board filtri (assignee, label va h.k.) — JSON
    filter       JSONB        NOT NULL DEFAULT '{}'::jsonb,
    created_by   UUID         NOT NULL REFERENCES users(id),
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX idx_boards_project_id ON boards (project_id) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_boards_updated_at
    BEFORE UPDATE ON boards
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- =============================================================================
-- BOARD_COLUMNS — board ustunlari (qaysi status'lar qaysi ustunda)
-- =============================================================================
CREATE TABLE board_columns (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    board_id     UUID         NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    name         VARCHAR(100) NOT NULL,
    position     INT          NOT NULL DEFAULT 0,
    wip_limit    INT          CHECK (wip_limit IS NULL OR wip_limit > 0),  -- WIP cheklov
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_board_columns_board_id ON board_columns (board_id, position);

-- =============================================================================
-- BOARD_COLUMN_STATUSES — bitta ustun bir nechta statusni o'z ichiga olishi mumkin
-- =============================================================================
CREATE TABLE board_column_statuses (
    column_id  UUID NOT NULL REFERENCES board_columns(id)    ON DELETE CASCADE,
    status_id  UUID NOT NULL REFERENCES workflow_statuses(id) ON DELETE CASCADE,
    PRIMARY KEY (column_id, status_id)
);

-- =============================================================================
-- AUDIT_LOGS — muhim amallar (login, ruxsat o'zgarishi, o'chirish)
-- =============================================================================
CREATE TABLE audit_logs (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID         REFERENCES users(id) ON DELETE SET NULL,
    action       VARCHAR(64)  NOT NULL,
    -- Misol: 'user.login', 'user.password_changed', 'project.deleted',
    --        'project.member_added', 'issue.deleted'

    entity_type  VARCHAR(32),
    entity_id    UUID,

    details      JSONB        NOT NULL DEFAULT '{}'::jsonb,
    ip_address   INET,
    user_agent   TEXT,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_user_id     ON audit_logs (user_id, created_at DESC);
CREATE INDEX idx_audit_logs_action      ON audit_logs (action, created_at DESC);
CREATE INDEX idx_audit_logs_entity      ON audit_logs (entity_type, entity_id) WHERE entity_id IS NOT NULL;
CREATE INDEX idx_audit_logs_created_at  ON audit_logs (created_at DESC);
