package issue_link

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, link *entity.IssueLink) (*entity.IssueLink, error)
	GetByID(ctx context.Context, id string) (*entity.IssueLink, error)
	ListByIssue(ctx context.Context, issueID string) ([]*entity.IssueLink, error)
	Delete(ctx context.Context, id, actorID string) error
}
