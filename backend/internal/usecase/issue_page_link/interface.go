package issue_page_link

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Link(ctx context.Context, issueID, pageID, actorID string) (*entity.IssuePageLink, error)
	Unlink(ctx context.Context, issueID, pageID string) error
	ListByIssue(ctx context.Context, issueID string) ([]*entity.IssuePageLink, error)
	ListByPage(ctx context.Context, pageID string) ([]*entity.IssuePageLink, error)
}
