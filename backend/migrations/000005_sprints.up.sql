-- =============================================================================
-- SPRINTS — sprint'lar
-- =============================================================================
CREATE TABLE sprints (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id   UUID         NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name         VARCHAR(255) NOT NULL,
    goal         TEXT,
    status       VARCHAR(32)  NOT NULL DEFAULT 'planned'
                 CHECK (status IN ('planned', 'active', 'completed')),
    start_date   DATE,
    end_date     DATE,
    started_at   TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_by   UUID         NOT NULL REFERENCES users(id),
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ,

    CHECK (end_date IS NULL OR start_date IS NULL OR end_date >= start_date)
);

CREATE INDEX idx_sprints_project_id ON sprints (project_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_sprints_status     ON sprints (project_id, status) WHERE deleted_at IS NULL;

-- Bir loyihada faqat bitta active sprint
CREATE UNIQUE INDEX uq_sprints_one_active_per_project
    ON sprints (project_id) WHERE status = 'active' AND deleted_at IS NULL;

CREATE TRIGGER trg_sprints_updated_at
    BEFORE UPDATE ON sprints
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
