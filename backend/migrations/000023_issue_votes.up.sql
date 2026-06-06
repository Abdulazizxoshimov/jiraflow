CREATE TABLE issue_votes (
    issue_id   UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES users(id)  ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (issue_id, user_id)
);

CREATE INDEX idx_issue_votes_issue ON issue_votes (issue_id);

ALTER TABLE issues ADD COLUMN vote_count INT NOT NULL DEFAULT 0;
