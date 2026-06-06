DROP TABLE IF EXISTS issue_worklogs CASCADE;
ALTER TABLE issues
    DROP COLUMN IF EXISTS original_estimate,
    DROP COLUMN IF EXISTS remaining_estimate;
