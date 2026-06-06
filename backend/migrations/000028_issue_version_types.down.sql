ALTER TABLE issue_versions DROP CONSTRAINT issue_versions_pkey;
ALTER TABLE issue_versions DROP COLUMN IF EXISTS version_type;
ALTER TABLE issue_versions ADD PRIMARY KEY (issue_id, version_id);
