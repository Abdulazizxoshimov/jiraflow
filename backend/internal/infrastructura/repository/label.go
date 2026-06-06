package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type LabelRepository interface {
	Create(ctx context.Context, l *entity.Label) error
	GetByID(ctx context.Context, id string) (*entity.Label, error)
	ListByProject(ctx context.Context, projectID string) ([]*entity.Label, error)
	Update(ctx context.Context, l *entity.Label) error
	Delete(ctx context.Context, id string) error
	GetByIDs(ctx context.Context, ids []string) ([]*entity.Label, error)
}
