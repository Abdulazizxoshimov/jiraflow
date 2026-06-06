-- =============================================================================
-- FAZA 3 — Confluence To'ldirish
-- =============================================================================

-- PAGE TAGS ──────────────────────────────────────────────────────────────────
CREATE TABLE page_tags (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    space_id   UUID        NOT NULL REFERENCES spaces(id) ON DELETE CASCADE,
    name       VARCHAR(64) NOT NULL,
    color      VARCHAR(16) NOT NULL DEFAULT '#6B7280',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (space_id, name)
);

CREATE TABLE page_tag_links (
    page_id    UUID NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    tag_id     UUID NOT NULL REFERENCES page_tags(id) ON DELETE CASCADE,
    PRIMARY KEY (page_id, tag_id)
);

CREATE INDEX idx_page_tag_links_tag_id ON page_tag_links (tag_id);

-- FAVORITES ──────────────────────────────────────────────────────────────────
CREATE TABLE user_favorites (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    entity_type VARCHAR(16) NOT NULL CHECK (entity_type IN ('page', 'space')),
    entity_id   UUID        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, entity_type, entity_id)
);

CREATE INDEX idx_user_favorites_user_id ON user_favorites (user_id, entity_type);

-- RECENTLY VISITED ───────────────────────────────────────────────────────────
CREATE TABLE page_views (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    page_id    UUID        NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    user_id    UUID        REFERENCES users(id) ON DELETE SET NULL,
    ip_address INET,
    viewed_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_page_views_page_id   ON page_views (page_id, viewed_at DESC);
CREATE INDEX idx_page_views_user_id   ON page_views (user_id, viewed_at DESC) WHERE user_id IS NOT NULL;
CREATE INDEX idx_page_views_analytics ON page_views (page_id, user_id, viewed_at DESC);

-- INLINE COMMENTS ────────────────────────────────────────────────────────────
CREATE TABLE inline_comments (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    page_id     UUID         NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    author_id   UUID         NOT NULL REFERENCES users(id),
    anchor_id   VARCHAR(128) NOT NULL,  -- TipTap node id yoki selection hash
    quote_text  TEXT,                   -- ta'kidlangan matn fragment
    body        TEXT         NOT NULL,
    resolved    BOOLEAN      NOT NULL DEFAULT FALSE,
    resolved_by UUID         REFERENCES users(id),
    resolved_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX idx_inline_comments_page_id ON inline_comments (page_id, anchor_id) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_inline_comments_updated_at
    BEFORE UPDATE ON inline_comments
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- PAGE TEMPLATES ─────────────────────────────────────────────────────────────
CREATE TABLE page_templates (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    space_id    UUID         REFERENCES spaces(id) ON DELETE CASCADE,  -- NULL = global
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    category    VARCHAR(64)  NOT NULL DEFAULT 'general',  -- general | meeting | retrospective | decision
    content     JSONB        NOT NULL DEFAULT '{}'::jsonb,
    content_text TEXT        NOT NULL DEFAULT '',
    icon        VARCHAR(64),
    created_by  UUID         NOT NULL REFERENCES users(id),
    is_global   BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_page_templates_space  ON page_templates (space_id) WHERE space_id IS NOT NULL;
CREATE INDEX idx_page_templates_global ON page_templates (is_global) WHERE is_global = TRUE;
CREATE INDEX idx_page_templates_cat    ON page_templates (category);

CREATE TRIGGER trg_page_templates_updated_at
    BEFORE UPDATE ON page_templates
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- PAGE RESTRICTIONS ──────────────────────────────────────────────────────────
CREATE TABLE page_restrictions (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    page_id     UUID        NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    type        VARCHAR(16) NOT NULL CHECK (type IN ('view', 'edit')),
    subject_type VARCHAR(16) NOT NULL CHECK (subject_type IN ('user', 'role')),
    subject_id  VARCHAR(128) NOT NULL,  -- user UUID yoki role nomi
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (page_id, type, subject_type, subject_id)
);

CREATE INDEX idx_page_restrictions_page ON page_restrictions (page_id, type);
