package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type IssuePageLinkRepository interface {
	Create(ctx context.Context, link *entity.IssuePageLink) error
	Delete(ctx context.Context, issueID, pageID string) error
	ListByIssue(ctx context.Context, issueID string) ([]*entity.IssuePageLink, error)
	ListByPage(ctx context.Context, pageID string) ([]*entity.IssuePageLink, error)
	Exists(ctx context.Context, issueID, pageID string) (bool, error)
}
