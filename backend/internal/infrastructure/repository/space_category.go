package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type SpaceCategoryRepository interface {
	Create(ctx context.Context, c *entity.SpaceCategory) error
	GetByID(ctx context.Context, id string) (*entity.SpaceCategory, error)
	List(ctx context.Context) ([]*entity.SpaceCategory, error)
	Update(ctx context.Context, c *entity.SpaceCategory) error
	Delete(ctx context.Context, id string) error
}
