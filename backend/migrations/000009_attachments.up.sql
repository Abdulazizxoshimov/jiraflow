-- =============================================================================
-- ATTACHMENTS — fayllar (issue, page, comment uchun)
-- =============================================================================
CREATE TABLE attachments (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_type   VARCHAR(16)  NOT NULL CHECK (parent_type IN ('issue', 'page', 'comment')),
    parent_id     UUID         NOT NULL,

    file_name     VARCHAR(255) NOT NULL,
    file_size     BIGINT       NOT NULL CHECK (file_size > 0),
    mime_type     VARCHAR(255) NOT NULL,
    storage_path  TEXT         NOT NULL,    -- S3 key yoki lokal path
    storage_type  VARCHAR(16)  NOT NULL DEFAULT 'local'
                  CHECK (storage_type IN ('local', 's3')),
    checksum      VARCHAR(64),              -- sha256 (dedup uchun ixtiyoriy)

    uploaded_by   UUID         NOT NULL REFERENCES users(id),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

CREATE INDEX idx_attachments_parent      ON attachments (parent_type, parent_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_attachments_uploaded_by ON attachments (uploaded_by) WHERE deleted_at IS NULL;
CREATE INDEX idx_attachments_created_at  ON attachments (created_at DESC) WHERE deleted_at IS NULL;

COMMENT ON COLUMN attachments.storage_path IS 'Lokal: relative path; S3: object key';
