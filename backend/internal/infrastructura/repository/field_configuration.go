package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type FieldConfigurationRepository interface {
	Create(ctx context.Context, c *entity.FieldConfiguration) error
	GetByID(ctx context.Context, id string) (*entity.FieldConfiguration, error)
	List(ctx context.Context, projectID *string) ([]*entity.FieldConfiguration, error)
	Delete(ctx context.Context, id string) error
}
