CREATE TABLE space_exports (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    space_id    UUID         NOT NULL REFERENCES spaces(id) ON DELETE CASCADE,
    requested_by UUID        NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    status      VARCHAR(20)  NOT NULL DEFAULT 'pending',  -- pending | processing | done | failed
    file_url    TEXT,
    error_msg   TEXT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_space_exports_space_id ON space_exports (space_id, created_at DESC);
