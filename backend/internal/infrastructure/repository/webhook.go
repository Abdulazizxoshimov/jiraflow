package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type WebhookRepository interface {
	Create(ctx context.Context, wh *entity.Webhook) error
	GetByID(ctx context.Context, id string) (*entity.Webhook, error)
	ListByProject(ctx context.Context, projectID string) ([]*entity.Webhook, error)
	ListBySpace(ctx context.Context, spaceID string) ([]*entity.Webhook, error)
	Update(ctx context.Context, wh *entity.Webhook) error
	Delete(ctx context.Context, id string) error
	FindByEvent(ctx context.Context, event string, projectID, spaceID *string) ([]*entity.Webhook, error)
	SaveDelivery(ctx context.Context, d *entity.WebhookDelivery) error
	ListDeliveries(ctx context.Context, webhookID string, limit int) ([]*entity.WebhookDelivery, error)
}
