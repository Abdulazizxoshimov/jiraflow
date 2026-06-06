CREATE TABLE issue_types (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    icon_url    TEXT,
    color       VARCHAR(20),
    is_subtask  BOOLEAN      NOT NULL DEFAULT false,
    is_system   BOOLEAN      NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE issue_type_schemes (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    project_id  UUID         REFERENCES projects(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE issue_type_scheme_members (
    scheme_id      UUID    NOT NULL REFERENCES issue_type_schemes(id) ON DELETE CASCADE,
    issue_type_id  UUID    NOT NULL REFERENCES issue_types(id) ON DELETE CASCADE,
    display_order  INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (scheme_id, issue_type_id)
);

-- Seed default system issue types
INSERT INTO issue_types (id, name, description, is_system) VALUES
  (gen_random_uuid(), 'Task',    'A task that needs to be done', true),
  (gen_random_uuid(), 'Bug',     'A problem which impairs functionality', true),
  (gen_random_uuid(), 'Story',   'A user story or feature request', true),
  (gen_random_uuid(), 'Epic',    'A large body of work', true),
  (gen_random_uuid(), 'Subtask', 'A subtask of another issue', true);

CREATE INDEX idx_issue_type_schemes_project ON issue_type_schemes (project_id);
