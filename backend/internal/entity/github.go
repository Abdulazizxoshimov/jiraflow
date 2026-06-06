package entity

import "time"

type GitHubRepo struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"project_id"`
	RepoFullName  string    `json:"repo_full_name"`
	RepoURL       string    `json:"repo_url"`
	WebhookSecret string    `json:"-"`
	CreatedBy     string    `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
}

type IssueCommit struct {
	ID          string     `json:"id"`
	IssueID     string     `json:"issue_id"`
	RepoID      string     `json:"repo_id"`
	SHA         string     `json:"sha"`
	Message     string     `json:"message"`
	AuthorName  string     `json:"author_name"`
	AuthorEmail string     `json:"author_email"`
	CommittedAt *time.Time `json:"committed_at,omitempty"`
	URL         string     `json:"url"`
	CreatedAt   time.Time  `json:"created_at"`
}

type IssuePullRequest struct {
	ID          string     `json:"id"`
	IssueID     string     `json:"issue_id"`
	RepoID      string     `json:"repo_id"`
	PRNumber    int        `json:"pr_number"`
	Title       string     `json:"title"`
	State       string     `json:"state"`
	URL         string     `json:"url"`
	AuthorLogin string     `json:"author_login"`
	CreatedAt   time.Time  `json:"created_at"`
	MergedAt    *time.Time `json:"merged_at,omitempty"`
}

type ConnectRepoReq struct {
	RepoFullName  string `json:"repo_full_name"  validate:"required"`
	RepoURL       string `json:"repo_url"        validate:"required,url"`
	WebhookSecret string `json:"webhook_secret"  validate:"required"`
}
