package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type GitHubRepoRepository interface {
	Create(ctx context.Context, repo *entity.GitHubRepo) error
	GetByProjectID(ctx context.Context, projectID string) (*entity.GitHubRepo, error)
	GetByRepoFullName(ctx context.Context, repoFullName string) (*entity.GitHubRepo, error)
	Delete(ctx context.Context, projectID string) error
}

type IssueCommitRepository interface {
	Create(ctx context.Context, commit *entity.IssueCommit) error
	ListByIssue(ctx context.Context, issueID string) ([]*entity.IssueCommit, error)
	GetBySHA(ctx context.Context, sha string) (*entity.IssueCommit, error)
}

type IssuePRRepository interface {
	Create(ctx context.Context, pr *entity.IssuePullRequest) error
	ListByIssue(ctx context.Context, issueID string) ([]*entity.IssuePullRequest, error)
	UpdateState(ctx context.Context, issueID string, prNumber int, state string, mergedAt *string) error
}
