ALTER TABLE boards
    ADD COLUMN IF NOT EXISTS swimlane_type VARCHAR(32) NOT NULL DEFAULT 'none';
-- swimlane_type: none | assignee | epic | priority | label
