CREATE TABLE blueprints (
    id             UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name           VARCHAR(255) NOT NULL,
    description    TEXT,
    icon_url       TEXT,
    category       VARCHAR(100),
    template_body  TEXT,
    schema         JSONB,
    is_system      BOOLEAN      NOT NULL DEFAULT false,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

INSERT INTO blueprints (id, name, description, category, template_body, is_system) VALUES
  (gen_random_uuid(), 'Meeting Notes', 'Capture meeting agenda and action items', 'meetings',
   E'## Attendees\n\n## Agenda\n\n## Discussion\n\n## Action Items\n\n| Action | Owner | Due Date |\n|--------|-------|----------|\n| | | |\n', true),
  (gen_random_uuid(), 'Decision', 'Document a decision and its rationale', 'decisions',
   E'## Background\n\n## Options Considered\n\n## Decision\n\n## Rationale\n\n## Consequences\n', true),
  (gen_random_uuid(), 'Retrospective', 'Sprint retrospective template', 'agile',
   E'## What went well\n\n## What could be improved\n\n## Action items\n\n| Action | Owner | Due Date |\n|--------|-------|----------|\n| | | |\n', true),
  (gen_random_uuid(), 'How-to Article', 'Step-by-step guide', 'documentation',
   E'## Overview\n\n## Prerequisites\n\n## Steps\n\n1. \n2. \n3. \n\n## Troubleshooting\n', true),
  (gen_random_uuid(), 'Product Requirements', 'Document product requirements and specs', 'product',
   E'## Problem Statement\n\n## Goals\n\n## Non-goals\n\n## Requirements\n\n## Open Questions\n', true);
