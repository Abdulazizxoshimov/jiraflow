package project

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, p *entity.Project, actorID string) (*entity.Project, error)
	GetByID(ctx context.Context, id string) (*entity.Project, error)
	GetByKey(ctx context.Context, key string) (*entity.Project, error)
	List(ctx context.Context, filter *entity.ProjectFilter) ([]*entity.Project, int, error)
	Update(ctx context.Context, id string, p *entity.Project, actorID string) (*entity.Project, error)
	Archive(ctx context.Context, id string, actorID string) error
	Delete(ctx context.Context, id string, actorID string) error
	GetLinkedSpace(ctx context.Context, projectID string) (*entity.Space, error)
	GetDashboard(ctx context.Context, projectID string) (*entity.ProjectDashboard, error)
}
