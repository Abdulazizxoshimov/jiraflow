package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type ProjectRepository interface {
	Create(ctx context.Context, p *entity.Project) error
	GetByID(ctx context.Context, id string) (*entity.Project, error)
	GetByKey(ctx context.Context, key string) (*entity.Project, error)
	List(ctx context.Context, filter *entity.ProjectFilter) ([]*entity.Project, int, error)
	Update(ctx context.Context, p *entity.Project) error
	SoftDelete(ctx context.Context, id string) error
	ExistsByKey(ctx context.Context, key string) (bool, error)
	IncrementIssueCounter(ctx context.Context, id string) (int64, error)
	// AllocateIssueCounters increments issue_counter by n atomically and returns
	// the first allocated number (start). Issues should use start, start+1, ..., start+n-1.
	AllocateIssueCounters(ctx context.Context, id string, n int) (start int64, err error)
	GetDashboard(ctx context.Context, projectID string) (*entity.ProjectDashboard, error)
}
