-- oauth_states: short-lived CSRF state tokens for OAuth2 flow
CREATE TABLE IF NOT EXISTS oauth_states (
    state       VARCHAR(255) PRIMARY KEY,
    redirect_url TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at  TIMESTAMPTZ NOT NULL DEFAULT (now() + INTERVAL '10 minutes')
);

-- oauth_accounts: social login accounts linked to users
CREATE TABLE IF NOT EXISTS oauth_accounts (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider         VARCHAR(50) NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    email            VARCHAR(255),
    name             VARCHAR(255),
    avatar_url       TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (provider, provider_user_id),
    UNIQUE (user_id, provider)
);

-- api_keys: service account tokens for CI/CD and integrations
CREATE TABLE IF NOT EXISTS api_keys (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name         VARCHAR(255) NOT NULL,
    key_prefix   VARCHAR(10) NOT NULL,
    key_hash     VARCHAR(255) NOT NULL UNIQUE,
    scopes       TEXT[]      NOT NULL DEFAULT '{}',
    last_used_at TIMESTAMPTZ,
    expires_at   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    revoked_at   TIMESTAMPTZ
);

-- permission_schemes: named sets of permissions assignable to projects
CREATE TABLE IF NOT EXISTS permission_schemes (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_by  UUID        REFERENCES users(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- permission_scheme_grants: individual permission grants inside a scheme
CREATE TABLE IF NOT EXISTS permission_scheme_grants (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    scheme_id   UUID        NOT NULL REFERENCES permission_schemes(id) ON DELETE CASCADE,
    permission  VARCHAR(100) NOT NULL,
    holder_type VARCHAR(50) NOT NULL, -- 'user' | 'role' | 'anyone'
    holder_id   UUID,                 -- NULL when holder_type = 'anyone'
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- partial unique indexes to handle NULL holder_id correctly
CREATE UNIQUE INDEX IF NOT EXISTS uq_scheme_grants_with_holder
    ON permission_scheme_grants (scheme_id, permission, holder_type, holder_id)
    WHERE holder_id IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS uq_scheme_grants_no_holder
    ON permission_scheme_grants (scheme_id, permission, holder_type)
    WHERE holder_id IS NULL;

ALTER TABLE projects ADD COLUMN IF NOT EXISTS permission_scheme_id UUID REFERENCES permission_schemes(id);

-- data_imports: track async import jobs (Jira XML, Trello JSON, Linear CSV)
CREATE TABLE IF NOT EXISTS data_imports (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID        NOT NULL REFERENCES users(id),
    source          VARCHAR(50) NOT NULL,
    status          VARCHAR(50) NOT NULL DEFAULT 'pending',
    total_items     INT         NOT NULL DEFAULT 0,
    processed_items INT         NOT NULL DEFAULT 0,
    error_message   TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at    TIMESTAMPTZ
);
