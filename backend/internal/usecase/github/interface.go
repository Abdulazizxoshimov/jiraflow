package github

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	ghinfra "github.com/jira-backend/jiraflow-backend/internal/infrastructure/github"
)

type UseCase interface {
	ConnectRepo(ctx context.Context, projectID, userID string, req *entity.ConnectRepoReq) (*entity.GitHubRepo, error)
	DisconnectRepo(ctx context.Context, projectID string) error
	GetRepo(ctx context.Context, projectID string) (*entity.GitHubRepo, error)
	HandlePushEvent(ctx context.Context, repoFullName string, commits []ghinfra.Commit) error
	HandlePREvent(ctx context.Context, repoFullName string, pr *ghinfra.PREvent) error
	ListCommits(ctx context.Context, issueID string) ([]*entity.IssueCommit, error)
	ListPRs(ctx context.Context, issueID string) ([]*entity.IssuePullRequest, error)
}
