-- =============================================================================
-- TIME TRACKING — original/remaining estimate + worklog
-- =============================================================================

-- Issues jadvaliga vaqt hisoblash ustunlari
ALTER TABLE issues
    ADD COLUMN IF NOT EXISTS original_estimate  INT,  -- sekundda
    ADD COLUMN IF NOT EXISTS remaining_estimate INT;  -- sekundda

-- Worklog — kim qancha vaqt sarfladi
CREATE TABLE issue_worklogs (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    issue_id     UUID         NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    user_id      UUID         NOT NULL REFERENCES users(id)  ON DELETE CASCADE,
    time_spent   INT          NOT NULL CHECK (time_spent > 0),  -- sekundda
    started_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    description  TEXT,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_worklogs_issue_id ON issue_worklogs (issue_id, created_at DESC);
CREATE INDEX idx_worklogs_user_id  ON issue_worklogs (user_id);

CREATE TRIGGER trg_worklogs_updated_at
    BEFORE UPDATE ON issue_worklogs
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
