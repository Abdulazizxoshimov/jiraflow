package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type IssueAssigneeRepository interface {
	Set(ctx context.Context, issueID string, userIDs []string, primaryID *string) error
	List(ctx context.Context, issueID string) ([]*entity.IssueAssignee, error)
	Remove(ctx context.Context, issueID, userID string) error
}
