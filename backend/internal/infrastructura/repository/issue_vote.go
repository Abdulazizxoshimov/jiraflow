package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type IssueVoteRepository interface {
	Toggle(ctx context.Context, issueID, userID string) (added bool, err error)
	GetSummary(ctx context.Context, issueID, viewerUserID string) (*entity.IssueVoteSummary, error)
}
