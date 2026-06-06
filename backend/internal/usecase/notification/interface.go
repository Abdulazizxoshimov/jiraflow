package notification

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Notify(ctx context.Context, n *entity.Notification) error
	ListByUser(ctx context.Context, userID string, filter *entity.NotificationFilter) ([]*entity.Notification, int, error)
	MarkRead(ctx context.Context, userID string, ids []string) error
	MarkAllRead(ctx context.Context, userID string) error
	Delete(ctx context.Context, id, userID string) error
	CountUnread(ctx context.Context, userID string) (int, error)
	GetPreference(ctx context.Context, userID string) (*entity.NotificationPreference, error)
	UpdatePreference(ctx context.Context, userID string, req *entity.UpdateNotificationPreferenceReq) (*entity.NotificationPreference, error)
}
