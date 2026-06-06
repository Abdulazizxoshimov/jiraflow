package page_lock

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Acquire(ctx context.Context, pageID, userID string, ttlSeconds int) (*entity.PageLock, error)
	Release(ctx context.Context, pageID, userID string) error
	Get(ctx context.Context, pageID string) (*entity.PageLock, error)
	Extend(ctx context.Context, pageID, userID string, ttlSeconds int) (*entity.PageLock, error)
}
