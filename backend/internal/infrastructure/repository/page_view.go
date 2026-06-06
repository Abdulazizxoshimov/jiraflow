package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type PageViewRepository interface {
	Record(ctx context.Context, view *entity.PageView) error
	GetAnalytics(ctx context.Context, pageID string) (*entity.PageAnalytics, error)
	ListRecentByUser(ctx context.Context, userID string, limit int) ([]*entity.RecentPage, error)
}
