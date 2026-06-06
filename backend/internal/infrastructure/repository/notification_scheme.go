package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type NotificationSchemeRepository interface {
	Create(ctx context.Context, s *entity.NotificationScheme) error
	GetByID(ctx context.Context, id string) (*entity.NotificationScheme, error)
	List(ctx context.Context) ([]*entity.NotificationScheme, error)
	Delete(ctx context.Context, id string) error
}
