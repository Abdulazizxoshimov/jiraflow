-- =============================================================================
-- COMMENTS — izohlar (polymorphic: issue yoki page uchun)
-- =============================================================================
CREATE TABLE comments (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_type  VARCHAR(16)  NOT NULL CHECK (parent_type IN ('issue', 'page')),
    parent_id    UUID         NOT NULL,  -- FK constraint kod tomonidan tekshiriladi
    author_id    UUID         NOT NULL REFERENCES users(id),
    content      JSONB        NOT NULL DEFAULT '{}'::jsonb,  -- TipTap JSON
    content_text TEXT         NOT NULL DEFAULT '',           -- plain (mention, notification uchun)

    -- Reply uchun (treaded comments)
    reply_to_id  UUID         REFERENCES comments(id) ON DELETE SET NULL,

    is_edited    BOOLEAN      NOT NULL DEFAULT FALSE,
    edited_at    TIMESTAMPTZ,

    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX idx_comments_parent       ON comments (parent_type, parent_id, created_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_comments_author       ON comments (author_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_comments_reply_to     ON comments (reply_to_id) WHERE deleted_at IS NULL AND reply_to_id IS NOT NULL;

CREATE TRIGGER trg_comments_updated_at
    BEFORE UPDATE ON comments
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMENT ON COLUMN comments.parent_type IS 'issue yoki page';
COMMENT ON COLUMN comments.parent_id IS 'FK ga issues.id yoki pages.id — kod tomonidan validatsiya';

-- =============================================================================
-- COMMENT_MENTIONS — izohdagi @mention'lar (notification uchun)
-- =============================================================================
CREATE TABLE comment_mentions (
    comment_id  UUID NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
    user_id     UUID NOT NULL REFERENCES users(id)    ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (comment_id, user_id)
);

CREATE INDEX idx_comment_mentions_user_id ON comment_mentions (user_id, created_at DESC);
