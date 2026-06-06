package telegram

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	GenerateCode(ctx context.Context, userID string) (code string, err error)
	GetStatus(ctx context.Context, userID string) (*entity.TelegramConnection, error)
	Disconnect(ctx context.Context, userID string) error
	HandleUpdate(ctx context.Context, update *entity.TelegramUpdate) error
	SendNotification(ctx context.Context, userID, message string) error
	SetupWebhook(ctx context.Context) error
}
