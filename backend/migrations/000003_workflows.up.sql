-- =============================================================================
-- WORKFLOWS — vazifalar uchun holatlar sxemasi
-- =============================================================================
CREATE TABLE workflows (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(255) NOT NULL,
    description   TEXT,
    is_default    BOOLEAN      NOT NULL DEFAULT FALSE,
    created_by    UUID         NOT NULL REFERENCES users(id),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

-- Faqat bitta default workflow bo'lishi mumkin
CREATE UNIQUE INDEX uq_workflows_default
    ON workflows (is_default) WHERE is_default = TRUE AND deleted_at IS NULL;

CREATE TRIGGER trg_workflows_updated_at
    BEFORE UPDATE ON workflows
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- =============================================================================
-- WORKFLOW_STATUSES — workflow ichidagi holatlar
-- =============================================================================
CREATE TABLE workflow_statuses (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id  UUID         NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    name         VARCHAR(100) NOT NULL,
    category     VARCHAR(32)  NOT NULL
                 CHECK (category IN ('todo', 'in_progress', 'done')),
    color        VARCHAR(7)   NOT NULL DEFAULT '#6B7280',  -- hex
    position     INT          NOT NULL DEFAULT 0,
    is_initial   BOOLEAN      NOT NULL DEFAULT FALSE,  -- yangi issue uchun standart status
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    UNIQUE (workflow_id, name)
);

CREATE INDEX idx_workflow_statuses_workflow_id ON workflow_statuses (workflow_id, position);

-- Har bir workflow uchun aniq bitta is_initial
CREATE UNIQUE INDEX uq_workflow_statuses_initial
    ON workflow_statuses (workflow_id) WHERE is_initial = TRUE;

CREATE TRIGGER trg_workflow_statuses_updated_at
    BEFORE UPDATE ON workflow_statuses
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- =============================================================================
-- WORKFLOW_TRANSITIONS — qaysi holat dan qaysi holatga o'tish mumkin
-- =============================================================================
CREATE TABLE workflow_transitions (
    id              UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id     UUID         NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    from_status_id  UUID         REFERENCES workflow_statuses(id) ON DELETE CASCADE,  -- NULL = istalgan holatdan
    to_status_id    UUID         NOT NULL REFERENCES workflow_statuses(id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,  -- "Start progress", "Resolve"
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    UNIQUE (workflow_id, from_status_id, to_status_id)
);

CREATE INDEX idx_workflow_transitions_workflow_id ON workflow_transitions (workflow_id);

COMMENT ON COLUMN workflow_transitions.from_status_id IS 'NULL = istalgan holatdan o''tish mumkin (global transition)';
