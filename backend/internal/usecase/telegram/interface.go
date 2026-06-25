package telegram

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	// GetStatus returns the Telegram connection for the given web user, or nil if not connected.
	GetStatus(ctx context.Context, userID string) (*entity.TelegramConnection, error)
	// VerifyCode links the web user's account using a code they received from the bot.
	VerifyCode(ctx context.Context, userID, code string) error
	Disconnect(ctx context.Context, userID string) error
	HandleUpdate(ctx context.Context, update *entity.TelegramUpdate) error
	SendNotification(ctx context.Context, userID, message string) error
	SetupWebhook(ctx context.Context) error
	DeleteWebhook(ctx context.Context) error
	BotInfo(ctx context.Context) *entity.TelegramBotInfo
	// StartPolling long-polls Telegram for updates (local dev, no public URL needed).
	StartPolling(ctx context.Context)
}
