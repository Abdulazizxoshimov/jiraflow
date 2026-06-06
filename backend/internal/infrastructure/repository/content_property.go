package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type ContentPropertyRepository interface {
	Set(ctx context.Context, p *entity.ContentProperty) error
	Get(ctx context.Context, entityType, entityID, key string) (*entity.ContentProperty, error)
	List(ctx context.Context, entityType, entityID string) ([]*entity.ContentProperty, error)
	Delete(ctx context.Context, entityType, entityID, key string) error
}
