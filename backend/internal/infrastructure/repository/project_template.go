package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type ProjectTemplateRepository interface {
	List(ctx context.Context) ([]*entity.ProjectTemplate, error)
	GetByID(ctx context.Context, id string) (*entity.ProjectTemplate, error)
}
