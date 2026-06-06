CREATE TABLE telegram_connections (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    telegram_id         BIGINT,
    chat_id             BIGINT,
    username            VARCHAR(255),
    verification_code   VARCHAR(8) UNIQUE,
    verified_at         TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tg_connections_user ON telegram_connections(user_id);
CREATE INDEX idx_tg_connections_chat ON telegram_connections(chat_id);
CREATE INDEX idx_tg_connections_code ON telegram_connections(verification_code);
