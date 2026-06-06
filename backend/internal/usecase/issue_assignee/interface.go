package issue_assignee

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Set(ctx context.Context, issueID string, req *entity.SetIssueAssigneesReq) error
	List(ctx context.Context, issueID string) ([]*entity.IssueAssignee, error)
	Remove(ctx context.Context, issueID, userID string) error
}
