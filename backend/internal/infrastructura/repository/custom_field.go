package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type CustomFieldRepository interface {
	Create(ctx context.Context, cf *entity.CustomField) error
	GetByID(ctx context.Context, id string) (*entity.CustomField, error)
	GetByKey(ctx context.Context, projectID, fieldKey string) (*entity.CustomField, error)
	ListByProject(ctx context.Context, projectID string) ([]*entity.CustomField, error)
	Update(ctx context.Context, cf *entity.CustomField) error
	Delete(ctx context.Context, id string) error
	ReorderFields(ctx context.Context, projectID string, positions map[string]int) error
}
