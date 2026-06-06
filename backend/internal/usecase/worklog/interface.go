package worklog

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, issueID, userID string, req *entity.CreateWorklogReq) (*entity.Worklog, error)
	GetByID(ctx context.Context, id string) (*entity.Worklog, error)
	List(ctx context.Context, filter *entity.WorklogFilter) ([]*entity.Worklog, int, error)
	Update(ctx context.Context, id, actorID string, req *entity.UpdateWorklogReq) (*entity.Worklog, error)
	Delete(ctx context.Context, id, actorID string) error
	GetTimeSummary(ctx context.Context, issueID string) (*entity.TimeSpentSummary, error)
	UpdateEstimates(ctx context.Context, issueID string, original, remaining *int) error
}
