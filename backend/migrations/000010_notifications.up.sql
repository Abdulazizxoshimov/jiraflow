-- =============================================================================
-- NOTIFICATIONS — bildirishnomalar
-- =============================================================================
CREATE TABLE notifications (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type          VARCHAR(64)  NOT NULL,
    -- Misol type lar:
    --   'issue_assigned', 'issue_status_changed', 'issue_commented',
    --   'mentioned', 'page_edited', 'page_commented'

    -- payload misoli:
    -- { "issue_id": "...", "issue_key": "PROJ-1", "actor_id": "...", "from_status": "...", "to_status": "..." }
    payload       JSONB        NOT NULL DEFAULT '{}'::jsonb,

    -- Bildirishnoma manbasi (deep-link uchun)
    entity_type   VARCHAR(16)  CHECK (entity_type IS NULL OR entity_type IN ('issue', 'page', 'comment')),
    entity_id     UUID,

    actor_id      UUID         REFERENCES users(id) ON DELETE SET NULL,  -- harakat qiluvchi

    read_at       TIMESTAMPTZ,
    email_sent_at TIMESTAMPTZ,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_user_unread
    ON notifications (user_id, created_at DESC) WHERE read_at IS NULL;
CREATE INDEX idx_notifications_user_all
    ON notifications (user_id, created_at DESC);
CREATE INDEX idx_notifications_entity
    ON notifications (entity_type, entity_id) WHERE entity_id IS NOT NULL;

-- =============================================================================
-- NOTIFICATION_PREFERENCES — foydalanuvchi sozlamalari
-- =============================================================================
CREATE TABLE notification_preferences (
    user_id            UUID         PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,

    -- Har bir trigger uchun email yuborilsinmi?
    email_assigned     BOOLEAN      NOT NULL DEFAULT TRUE,
    email_mentioned    BOOLEAN      NOT NULL DEFAULT TRUE,
    email_commented    BOOLEAN      NOT NULL DEFAULT TRUE,
    email_status       BOOLEAN      NOT NULL DEFAULT FALSE,
    email_watcher      BOOLEAN      NOT NULL DEFAULT FALSE,

    -- Daily digest
    daily_digest       BOOLEAN      NOT NULL DEFAULT FALSE,

    updated_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_notification_preferences_updated_at
    BEFORE UPDATE ON notification_preferences
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
