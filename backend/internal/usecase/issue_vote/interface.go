package issue_vote

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Toggle(ctx context.Context, issueID, actorID string) (added bool, err error)
	GetSummary(ctx context.Context, issueID, actorID string) (*entity.IssueVoteSummary, error)
}
