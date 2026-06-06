CREATE TABLE content_properties (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type  VARCHAR(50)  NOT NULL,  -- page | issue | space
    entity_id    UUID         NOT NULL,
    key          VARCHAR(255) NOT NULL,
    value        JSONB        NOT NULL DEFAULT '{}',
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (entity_type, entity_id, key)
);

CREATE INDEX idx_content_properties_entity ON content_properties (entity_type, entity_id);
