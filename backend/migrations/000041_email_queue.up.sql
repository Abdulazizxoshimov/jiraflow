CREATE TABLE IF NOT EXISTS email_queue (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "to"            TEXT[]          NOT NULL,
    subject         TEXT            NOT NULL,
    body            TEXT            NOT NULL,
    status          TEXT            NOT NULL DEFAULT 'pending', -- pending | sent | failed
    attempts        INT             NOT NULL DEFAULT 0,
    next_attempt_at TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    sent_at         TIMESTAMPTZ,
    error           TEXT,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_email_queue_pending
    ON email_queue(next_attempt_at)
    WHERE status = 'pending';
