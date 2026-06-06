package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type TelegramRepository interface {
	Create(ctx context.Context, conn *entity.TelegramConnection) error
	GetByUserID(ctx context.Context, userID string) (*entity.TelegramConnection, error)
	GetByChatID(ctx context.Context, chatID int64) (*entity.TelegramConnection, error)
	GetByVerificationCode(ctx context.Context, code string) (*entity.TelegramConnection, error)
	UpdateVerified(ctx context.Context, id string, telegramID, chatID int64, username string) error
	Delete(ctx context.Context, userID string) error
}
