package activity_feed

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Record(ctx context.Context, event *entity.ActivityEvent) error
	List(ctx context.Context, filter *entity.ActivityFilter) ([]*entity.ActivityEvent, int, error)
}
