-- Add version_type to distinguish fix versions from affected versions.
-- Existing rows default to 'fix' (backward compatible).
ALTER TABLE issue_versions
    ADD COLUMN version_type VARCHAR(20) NOT NULL DEFAULT 'fix'
        CHECK (version_type IN ('fix', 'affects'));

-- Recreate PK to allow same version as both fix and affects.
ALTER TABLE issue_versions DROP CONSTRAINT issue_versions_pkey;
ALTER TABLE issue_versions ADD PRIMARY KEY (issue_id, version_id, version_type);
