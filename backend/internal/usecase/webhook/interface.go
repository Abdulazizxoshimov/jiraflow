package webhook

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, actorID string, req *entity.CreateWebhookReq) (*entity.Webhook, error)
	GetByID(ctx context.Context, id string) (*entity.Webhook, error)
	ListByProject(ctx context.Context, projectID string) ([]*entity.Webhook, error)
	ListBySpace(ctx context.Context, spaceID string) ([]*entity.Webhook, error)
	Update(ctx context.Context, id string, req *entity.UpdateWebhookReq) (*entity.Webhook, error)
	Delete(ctx context.Context, id string) error
	Trigger(ctx context.Context, event string, projectID, spaceID *string, payload map[string]any) error
	ListDeliveries(ctx context.Context, webhookID string, limit int) ([]*entity.WebhookDelivery, error)
}
