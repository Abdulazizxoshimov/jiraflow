package workflow

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, wf *entity.Workflow) (*entity.Workflow, error)
	GetByID(ctx context.Context, id string) (*entity.Workflow, error)
	GetWithDetails(ctx context.Context, id string) (*entity.Workflow, error)
	List(ctx context.Context, filter *entity.Filter) ([]*entity.Workflow, int, error)
	Update(ctx context.Context, id string, wf *entity.Workflow) (*entity.Workflow, error)
	SetDefault(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
	CreateStatus(ctx context.Context, s *entity.WorkflowStatus) (*entity.WorkflowStatus, error)
	UpdateStatus(ctx context.Context, s *entity.WorkflowStatus) (*entity.WorkflowStatus, error)
	DeleteStatus(ctx context.Context, id string) error
	CreateTransition(ctx context.Context, t *entity.WorkflowTransition) (*entity.WorkflowTransition, error)
	DeleteTransition(ctx context.Context, id string) error
}
