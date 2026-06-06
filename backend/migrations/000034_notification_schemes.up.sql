CREATE TABLE notification_schemes (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE notification_scheme_rules (
    id              UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    scheme_id       UUID         NOT NULL REFERENCES notification_schemes(id) ON DELETE CASCADE,
    event_type      VARCHAR(100) NOT NULL,
    recipient_type  VARCHAR(50)  NOT NULL,  -- role | user | reporter | assignee | watchers
    recipient_id    UUID,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

ALTER TABLE projects ADD COLUMN IF NOT EXISTS notification_scheme_id UUID REFERENCES notification_schemes(id) ON DELETE SET NULL;

CREATE INDEX idx_notification_scheme_rules_scheme ON notification_scheme_rules (scheme_id);
