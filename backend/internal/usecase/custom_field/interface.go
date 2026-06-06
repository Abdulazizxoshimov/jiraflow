package custom_field

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, projectID string, cf *entity.CustomField) (*entity.CustomField, error)
	GetByID(ctx context.Context, id string) (*entity.CustomField, error)
	ListByProject(ctx context.Context, projectID string) ([]*entity.CustomField, error)
	Update(ctx context.Context, id string, cf *entity.CustomField) (*entity.CustomField, error)
	Delete(ctx context.Context, id string) error
	Reorder(ctx context.Context, projectID string, positions map[string]int) error
}
