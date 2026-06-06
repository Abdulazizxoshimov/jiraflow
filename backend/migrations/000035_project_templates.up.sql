CREATE TABLE project_templates (
    id                      UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name                    VARCHAR(255) NOT NULL,
    type                    VARCHAR(50)  NOT NULL,  -- scrum | kanban | business
    description             TEXT,
    icon_url                TEXT,
    default_workflow_config JSONB,
    is_system               BOOLEAN      NOT NULL DEFAULT false,
    created_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

INSERT INTO project_templates (id, name, type, description, is_system, default_workflow_config) VALUES
  (gen_random_uuid(), 'Scrum',   'scrum',   'Agile Scrum project with sprints', true,
   '{"statuses": ["To Do", "In Progress", "In Review", "Done"]}'),
  (gen_random_uuid(), 'Kanban',  'kanban',  'Kanban board for continuous flow', true,
   '{"statuses": ["Backlog", "In Progress", "Done"]}'),
  (gen_random_uuid(), 'Business','business','Simple project for business teams', true,
   '{"statuses": ["To Do", "In Progress", "Done"]}');
