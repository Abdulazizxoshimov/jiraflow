CREATE TABLE field_configurations (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(255) NOT NULL,
    project_id UUID         REFERENCES projects(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE field_config_items (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    config_id    UUID         NOT NULL REFERENCES field_configurations(id) ON DELETE CASCADE,
    field_name   VARCHAR(100) NOT NULL,
    is_required  BOOLEAN      NOT NULL DEFAULT false,
    is_hidden    BOOLEAN      NOT NULL DEFAULT false,
    description  TEXT,
    UNIQUE (config_id, field_name)
);

CREATE INDEX idx_field_configurations_project ON field_configurations (project_id);
