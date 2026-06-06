package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type WorklogRepository interface {
	Create(ctx context.Context, w *entity.Worklog) error
	GetByID(ctx context.Context, id string) (*entity.Worklog, error)
	List(ctx context.Context, filter *entity.WorklogFilter) ([]*entity.Worklog, int, error)
	Update(ctx context.Context, w *entity.Worklog) error
	Delete(ctx context.Context, id string) error
	SumByIssue(ctx context.Context, issueID string) (int, error)
}
