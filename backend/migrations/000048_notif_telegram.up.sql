ALTER TABLE notification_preferences
    ADD COLUMN IF NOT EXISTS telegram_enabled BOOLEAN NOT NULL DEFAULT TRUE;
