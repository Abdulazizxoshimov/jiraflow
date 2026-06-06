package repository

import (
	"context"
	"time"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type PageLockRepository interface {
	Acquire(ctx context.Context, lock *entity.PageLock) error
	Release(ctx context.Context, pageID, userID string) error
	Get(ctx context.Context, pageID string) (*entity.PageLock, error)
	Extend(ctx context.Context, pageID, userID string, expiresAt time.Time) error
	Cleanup(ctx context.Context) error
}
