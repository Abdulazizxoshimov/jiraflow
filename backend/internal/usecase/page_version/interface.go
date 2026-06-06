package page_version

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	GetByID(ctx context.Context, id string) (*entity.PageVersion, error)
	GetByVersion(ctx context.Context, pageID string, version int) (*entity.PageVersion, error)
	ListByPage(ctx context.Context, pageID string, filter *entity.Filter) ([]*entity.PageVersion, int, error)
	Diff(ctx context.Context, pageID string, v1, v2 int) (*entity.PageVersionDiff, error)
}
