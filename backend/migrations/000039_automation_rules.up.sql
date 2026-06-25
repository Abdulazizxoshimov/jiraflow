CREATE TABLE automation_rules (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id   UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name         TEXT NOT NULL,
    description  TEXT,
    trigger_type TEXT NOT NULL,
    trigger_config JSONB NOT NULL DEFAULT '{}',
    conditions   JSONB NOT NULL DEFAULT '[]',
    actions      JSONB NOT NULL DEFAULT '[]',
    is_active    BOOLEAN NOT NULL DEFAULT TRUE,
    created_by   UUID NOT NULL REFERENCES users(id),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_automation_rules_project ON automation_rules(project_id);
CREATE INDEX idx_automation_rules_trigger ON automation_rules(trigger_type) WHERE is_active = TRUE;

CREATE TABLE automation_logs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rule_id     UUID NOT NULL REFERENCES automation_rules(id) ON DELETE CASCADE,
    entity_id   TEXT,
    entity_type TEXT,
    status      TEXT NOT NULL DEFAULT 'success',
    executed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    error_msg   TEXT
);

CREATE INDEX idx_automation_logs_rule ON automation_logs(rule_id, executed_at DESC);
