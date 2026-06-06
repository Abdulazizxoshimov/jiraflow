package telegram

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	tgclient "github.com/jira-backend/jiraflow-backend/internal/infrastructure/telegram"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
)

type useCase struct {
	repo        repository.TelegramRepository
	tg          *tgclient.Client
	webhookURL  string
	webhookSecret string
}

func New(
	repo repository.TelegramRepository,
	tg *tgclient.Client,
	webhookURL, webhookSecret string,
) UseCase {
	return &useCase{
		repo:          repo,
		tg:            tg,
		webhookURL:    webhookURL,
		webhookSecret: webhookSecret,
	}
}

func (uc *useCase) GenerateCode(ctx context.Context, userID string) (string, error) {
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	conn := &entity.TelegramConnection{
		ID:               uuid.NewString(),
		UserID:           userID,
		VerificationCode: &code,
		CreatedAt:        time.Now().UTC(),
	}
	if err := uc.repo.Create(ctx, conn); err != nil {
		return "", fmt.Errorf("telegram.GenerateCode: %w", err)
	}
	return code, nil
}

func (uc *useCase) GetStatus(ctx context.Context, userID string) (*entity.TelegramConnection, error) {
	conn, err := uc.repo.GetByUserID(ctx, userID)
	if err != nil {
		if appErr, ok := err.(interface{ Code() string }); ok && appErr.Code() == "NOT_FOUND" {
			return nil, nil
		}
		return nil, err
	}
	return conn, nil
}

func (uc *useCase) Disconnect(ctx context.Context, userID string) error {
	return uc.repo.Delete(ctx, userID)
}

func (uc *useCase) HandleUpdate(ctx context.Context, update *entity.TelegramUpdate) error {
	if update.Message == nil {
		return nil
	}
	msg := update.Message
	text := strings.TrimSpace(msg.Text)
	chatID := msg.Chat.ID
	tgUser := msg.From

	switch {
	case strings.HasPrefix(text, "/start"):
		parts := strings.Fields(text)
		if len(parts) < 2 {
			_ = uc.tg.SendMessage(ctx, chatID, "Assalomu alaykum! JiraFlow'da /start <kod> yuboring.")
			return nil
		}
		code := parts[1]
		conn, err := uc.repo.GetByVerificationCode(ctx, code)
		if err != nil {
			_ = uc.tg.SendMessage(ctx, chatID, "Kod topilmadi yoki muddati o'tgan.")
			return nil
		}
		username := tgUser.Username
		if err := uc.repo.UpdateVerified(ctx, conn.ID, tgUser.ID, chatID, username); err != nil {
			return fmt.Errorf("telegram.HandleUpdate UpdateVerified: %w", err)
		}
		_ = uc.tg.SendMessage(ctx, chatID, fmt.Sprintf("✅ Tabriklaymiz! JiraFlow akkauntingiz ulandi, %s.", tgUser.FirstName))

	case text == "/myissues":
		_ = uc.tg.SendMessage(ctx, chatID, "Bu funksiya tez orada qo'shiladi.")

	case text == "/status":
		conn, err := uc.repo.GetByChatID(ctx, chatID)
		if err != nil {
			_ = uc.tg.SendMessage(ctx, chatID, "Siz hali JiraFlow akkauntingizni ulamamgansiz.")
			return nil
		}
		_ = uc.tg.SendMessage(ctx, chatID, fmt.Sprintf("✅ Ulangan. UserID: %s", conn.UserID))

	case text == "/help":
		_ = uc.tg.SendMessage(ctx, chatID, "/start <kod> - akkauntni ulash\n/status - ulash holatini ko'rish\n/help - yordam")

	default:
		_ = uc.tg.SendMessage(ctx, chatID, "Noma'lum buyruq. /help ni ko'ring.")
	}
	return nil
}

func (uc *useCase) SendNotification(ctx context.Context, userID, message string) error {
	conn, err := uc.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil // not connected — silent
	}
	if conn.ChatID == nil || conn.VerifiedAt == nil {
		return nil
	}
	return uc.tg.SendMessage(ctx, *conn.ChatID, message)
}

func (uc *useCase) SetupWebhook(ctx context.Context) error {
	if uc.webhookURL == "" {
		return apperr.BadRequest("TELEGRAM_WEBHOOK_URL not configured")
	}
	return uc.tg.SetWebhook(ctx, uc.webhookURL, uc.webhookSecret)
}
