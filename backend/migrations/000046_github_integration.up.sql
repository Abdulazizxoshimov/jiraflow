CREATE TABLE github_repos (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id      UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    repo_full_name  VARCHAR(255) NOT NULL,
    repo_url        TEXT NOT NULL,
    webhook_secret  VARCHAR(255) NOT NULL,
    created_by      UUID NOT NULL REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(project_id)
);

CREATE TABLE issue_commits (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    issue_id        UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    repo_id         UUID NOT NULL REFERENCES github_repos(id) ON DELETE CASCADE,
    sha             VARCHAR(40) NOT NULL,
    message         TEXT NOT NULL,
    author_name     VARCHAR(255),
    author_email    VARCHAR(255),
    committed_at    TIMESTAMPTZ,
    url             TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE issue_pull_requests (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    issue_id        UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    repo_id         UUID NOT NULL REFERENCES github_repos(id) ON DELETE CASCADE,
    pr_number       INT NOT NULL,
    title           TEXT NOT NULL,
    state           VARCHAR(20) NOT NULL DEFAULT 'open',
    url             TEXT NOT NULL,
    author_login    VARCHAR(255),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    merged_at       TIMESTAMPTZ
);

CREATE INDEX idx_issue_commits_issue ON issue_commits(issue_id);
CREATE INDEX idx_issue_commits_sha ON issue_commits(sha);
CREATE INDEX idx_issue_prs_issue ON issue_pull_requests(issue_id);
