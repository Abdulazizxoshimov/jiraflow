package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type PageTemplateRepository interface {
	Create(ctx context.Context, t *entity.PageTemplate) error
	GetByID(ctx context.Context, id string) (*entity.PageTemplate, error)
	List(ctx context.Context, filter *entity.PageTemplateFilter) ([]*entity.PageTemplate, int, error)
	Update(ctx context.Context, t *entity.PageTemplate) error
	Delete(ctx context.Context, id string) error
}
