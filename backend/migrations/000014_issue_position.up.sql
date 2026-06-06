-- Issue'larga backlog drag-and-drop tartibi uchun position ustuni
ALTER TABLE issues ADD COLUMN IF NOT EXISTS position INT NOT NULL DEFAULT 0;

CREATE INDEX IF NOT EXISTS idx_issues_position ON issues (project_id, position) WHERE deleted_at IS NULL;

-- Mavjud issue'larni issue_number tartibida pozitsiyalash
UPDATE issues SET position = issue_number WHERE position = 0;
