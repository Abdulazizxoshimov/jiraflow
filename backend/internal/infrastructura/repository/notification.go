package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type NotificationRepository interface {
	Create(ctx context.Context, n *entity.Notification) error
	GetByID(ctx context.Context, id string) (*entity.Notification, error)
	ListByUser(ctx context.Context, userID string, filter *entity.NotificationFilter) ([]*entity.Notification, int, error)
	MarkRead(ctx context.Context, userID string, ids []string) error
	MarkAllRead(ctx context.Context, userID string) error
	Delete(ctx context.Context, id string) error
	CountUnread(ctx context.Context, userID string) (int, error)

	GetPreference(ctx context.Context, userID string) (*entity.NotificationPreference, error)
	UpsertPreference(ctx context.Context, pref *entity.NotificationPreference) error
}
