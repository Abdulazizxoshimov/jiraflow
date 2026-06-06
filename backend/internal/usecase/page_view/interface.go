package page_view

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	RecordView(ctx context.Context, pageID string, userID *string, ipAddress *string) error
	GetAnalytics(ctx context.Context, pageID string) (*entity.PageAnalytics, error)
	ListRecentByUser(ctx context.Context, userID string, limit int) ([]*entity.RecentPage, error)
}
