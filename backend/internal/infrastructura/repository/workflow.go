package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type WorkflowRepository interface {
	Create(ctx context.Context, wf *entity.Workflow) error
	GetByID(ctx context.Context, id string) (*entity.Workflow, error)
	GetWithDetails(ctx context.Context, id string) (*entity.Workflow, error)
	List(ctx context.Context, filter *entity.Filter) ([]*entity.Workflow, int, error)
	Update(ctx context.Context, wf *entity.Workflow) error
	SoftDelete(ctx context.Context, id string) error
	SetDefault(ctx context.Context, id string) error
	GetDefault(ctx context.Context) (*entity.Workflow, error)

	CreateStatus(ctx context.Context, s *entity.WorkflowStatus) error
	GetStatusByID(ctx context.Context, id string) (*entity.WorkflowStatus, error)
	ListStatuses(ctx context.Context, workflowID string) ([]*entity.WorkflowStatus, error)
	UpdateStatus(ctx context.Context, s *entity.WorkflowStatus) error
	DeleteStatus(ctx context.Context, id string) error

	CreateTransition(ctx context.Context, t *entity.WorkflowTransition) error
	GetTransitionByID(ctx context.Context, id string) (*entity.WorkflowTransition, error)
	ListTransitions(ctx context.Context, workflowID string) ([]*entity.WorkflowTransition, error)
	DeleteTransition(ctx context.Context, id string) error
	IsTransitionAllowed(ctx context.Context, workflowID, fromStatusID, toStatusID string) (bool, error)
}
