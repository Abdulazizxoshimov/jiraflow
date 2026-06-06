package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type PageVersionRepository interface {
	Create(ctx context.Context, v *entity.PageVersion) error
	GetByID(ctx context.Context, id string) (*entity.PageVersion, error)
	GetByVersion(ctx context.Context, pageID string, version int) (*entity.PageVersion, error)
	ListByPage(ctx context.Context, pageID string, filter *entity.Filter) ([]*entity.PageVersion, int, error)
	GetLatestVersion(ctx context.Context, pageID string) (int, error)
}
