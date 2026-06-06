package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type BlueprintRepository interface {
	Create(ctx context.Context, b *entity.Blueprint) error
	GetByID(ctx context.Context, id string) (*entity.Blueprint, error)
	List(ctx context.Context) ([]*entity.Blueprint, error)
	Delete(ctx context.Context, id string) error
}
