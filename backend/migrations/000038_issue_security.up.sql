CREATE TABLE IF NOT EXISTS security_schemes (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    project_id  UUID REFERENCES projects(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS security_levels (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scheme_id   UUID NOT NULL REFERENCES security_schemes(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS security_level_members (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    level_id   UUID NOT NULL REFERENCES security_levels(id) ON DELETE CASCADE,
    type       VARCHAR(50) NOT NULL CHECK (type IN ('role','user','group','reporter','assignee')),
    value      VARCHAR(255)
);

ALTER TABLE issues ADD COLUMN IF NOT EXISTS security_level_id UUID REFERENCES security_levels(id) ON DELETE SET NULL;
