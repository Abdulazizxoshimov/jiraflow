package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/redis"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	tgclient "github.com/jira-backend/jiraflow-backend/internal/infrastructure/telegram"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
)

var issueKeyRe = regexp.MustCompile(`(?i)^[A-Z]{1,10}-\d+$`)

const pendingTTL = 10 * time.Minute

// pendingConn holds Telegram user info before the web user verifies the code.
type pendingConn struct {
	ChatID     int64  `json:"chat_id"`
	TelegramID int64  `json:"telegram_id"`
	Username   string `json:"username"`
	FirstName  string `json:"first_name"`
}

type useCase struct {
	repo          repository.TelegramRepository
	issueRepo     repository.IssueRepository
	cache         redis.Cache
	tg            *tgclient.Client
	webhookURL    string
	webhookSecret string
	botUsername   string
}

func New(
	repo repository.TelegramRepository,
	issueRepo repository.IssueRepository,
	cache redis.Cache,
	tg *tgclient.Client,
	webhookURL, webhookSecret, botUsername string,
) UseCase {
	return &useCase{
		repo:          repo,
		issueRepo:     issueRepo,
		cache:         cache,
		tg:            tg,
		webhookURL:    webhookURL,
		webhookSecret: webhookSecret,
		botUsername:   botUsername,
	}
}

func (uc *useCase) pendingKey(code string) string { return "tg:pending:" + code }

// storePending saves Telegram user info in Redis keyed by code.
// Pass struct directly — cache.Set already calls json.Marshal internally.
func (uc *useCase) storePending(ctx context.Context, code string, p pendingConn) error {
	if uc.cache == nil {
		return fmt.Errorf("cache not configured")
	}
	return uc.cache.Set(ctx, uc.pendingKey(code), p, pendingTTL)
}

// popPending retrieves and deletes a pending connection by code.
func (uc *useCase) popPending(ctx context.Context, code string) (*pendingConn, error) {
	if uc.cache == nil {
		return nil, fmt.Errorf("cache not configured")
	}
	raw, err := uc.cache.Get(ctx, uc.pendingKey(code))
	if err != nil || raw == "" {
		return nil, nil
	}
	_ = uc.cache.Del(ctx, uc.pendingKey(code))
	var p pendingConn
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func generateCode() string { return fmt.Sprintf("%06d", rand.Intn(1000000)) }

// ─── UseCase interface ────────────────────────────────────────────────────────

// GenerateCode is kept for backwards compat but no longer used in the primary flow.
func (uc *useCase) GenerateCode(ctx context.Context, userID string) (string, error) {
	return "", apperr.BadRequest("use the bot to get your connection code: send /start to @" + uc.botUsername)
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

// VerifyCode links a web user account using a code received from the bot.
func (uc *useCase) VerifyCode(ctx context.Context, userID, code string) error {
	p, err := uc.popPending(ctx, code)
	if err != nil {
		return fmt.Errorf("telegram.VerifyCode: %w", err)
	}
	if p == nil {
		return apperr.BadRequest("invalid or expired code — send /start to @" + uc.botUsername + " to get a new one")
	}
	if err := uc.repo.Link(ctx, userID, p.TelegramID, p.ChatID, p.Username); err != nil {
		return fmt.Errorf("telegram.VerifyCode link: %w", err)
	}
	// Notify the Telegram user that linking succeeded
	go func() {
		_ = uc.tg.SendMessage(context.Background(), p.ChatID, fmt.Sprintf(
			"✅ <b>Muvaffaqiyatli ulandi!</b>\n\nJiraFlow akkauntingiz bog'landi, %s.\n\n"+
				"📋 /myissues — vazifalaringizni ko'ring",
			p.FirstName,
		))
	}()
	return nil
}

// HandleUpdate processes an incoming Telegram message.
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
		// New flow: bot generates the code and sends it to the user.
		// The user then enters this code on the website.
		code := generateCode()
		pending := pendingConn{
			ChatID:     chatID,
			TelegramID: tgUser.ID,
			Username:   tgUser.Username,
			FirstName:  tgUser.FirstName,
		}
		if err := uc.storePending(ctx, code, pending); err != nil {
			_ = uc.tg.SendMessage(ctx, chatID, "❌ Xatolik yuz berdi. Iltimos, keyinroq urinib ko'ring.")
			return err
		}
		_ = uc.tg.SendMessage(ctx, chatID, fmt.Sprintf(
			"👋 Salom, <b>%s</b>!\n\n"+
				"JiraFlow akkauntingizni ulash uchun quyidagi kodni saytga kiriting:\n\n"+
				"🔑 <code>%s</code>\n\n"+
				"<b>Qayerga kiritish:</b>\n"+
				"JiraFlow → Settings → Integrations → Telegram → \"Kodni kiriting\"\n\n"+
				"⏳ Kod <b>10 daqiqa</b> amal qiladi.",
			tgUser.FirstName, code,
		))

	case text == "/myissues":
		uc.handleMyIssues(ctx, chatID)

	case text == "/status":
		conn, err := uc.repo.GetByChatID(ctx, chatID)
		if err != nil {
			_ = uc.tg.SendMessage(ctx, chatID,
				"❌ Siz hali JiraFlow akkauntingizni ulamamgansiz.\n\n"+
					"/start yuboring — kod olasiz va saytda kiriting.")
			return nil
		}
		username := ""
		if conn.Username != nil {
			username = "@" + *conn.Username
		}
		_ = uc.tg.SendMessage(ctx, chatID, fmt.Sprintf(
			"✅ <b>Ulangan</b>\nTelegram: %s\nJiraFlow UserID: <code>%s</code>",
			username, conn.UserID,
		))

	case text == "/help":
		uc.sendHelp(ctx, chatID)

	case issueKeyRe.MatchString(text):
		uc.handleIssueKey(ctx, chatID, strings.ToUpper(text))

	default:
		_ = uc.tg.SendMessage(ctx, chatID,
			"Noma'lum buyruq. /help ni ko'ring yoki issue key yuboring (masalan: <code>PROJ-123</code>).")
	}
	return nil
}

func (uc *useCase) handleMyIssues(ctx context.Context, chatID int64) {
	conn, err := uc.repo.GetByChatID(ctx, chatID)
	if err != nil {
		_ = uc.tg.SendMessage(ctx, chatID, "❌ Akkauntingiz ulanmagan.\n/start yuboring — kod olasiz.")
		return
	}

	issues, _, err := uc.issueRepo.List(ctx, &entity.IssueFilter{
		Filter:     entity.Filter{Page: 1, Limit: 10},
		AssigneeID: conn.UserID,
	})
	if err != nil || len(issues) == 0 {
		_ = uc.tg.SendMessage(ctx, chatID, "📭 Hozircha sizga biriktirilgan vazifalar yo'q.")
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📋 <b>Sizning vazifalaringiz</b> (%d ta):\n\n", len(issues)))
	for i, issue := range issues {
		status := "—"
		if issue.Status != nil {
			status = issue.Status.Name
		}
		sb.WriteString(fmt.Sprintf(
			"%d. <b>#%d</b> — %s\n   📌 %s\n\n",
			i+1, issue.IssueNumber, issue.Title, status,
		))
	}
	sb.WriteString("💡 Issue key yuboring (masalan: <code>PROJ-1</code>) — batafsil ma'lumot.")
	_ = uc.tg.SendMessage(ctx, chatID, sb.String())
}

func (uc *useCase) handleIssueKey(ctx context.Context, chatID int64, key string) {
	issue, err := uc.issueRepo.GetByKey(ctx, key)
	if err != nil {
		_ = uc.tg.SendMessage(ctx, chatID, fmt.Sprintf("❌ <code>%s</code> topilmadi.", key))
		return
	}
	status := "—"
	if issue.Status != nil {
		status = issue.Status.Name
	}
	priority := issue.Priority
	if priority == "" {
		priority = "—"
	}
	assignee := "Biriktirilmagan"
	if issue.Assignee != nil {
		assignee = issue.Assignee.FullName
	}
	_ = uc.tg.SendMessage(ctx, chatID, fmt.Sprintf(
		"🎫 <b>%s</b> — #%d\n\n📝 %s\n\n📌 Status: <b>%s</b>\n⚡️ Muhimlik: %s\n👤 Biriktirilgan: %s",
		key, issue.IssueNumber, issue.Title, status, priority, assignee,
	))
}

func (uc *useCase) sendHelp(ctx context.Context, chatID int64) {
	_ = uc.tg.SendMessage(ctx, chatID,
		"❓ <b>JiraFlow Bot — Yordam</b>\n\n"+
			"<b>Buyruqlar:</b>\n"+
			"/start — ulash kodi olish\n"+
			"/myissues — menga biriktirilgan vazifalar\n"+
			"/status — ulash holati\n"+
			"/help — ushbu yordam\n\n"+
			"<b>Issue qidirish:</b>\n"+
			"Issue key yuboring (masalan: <code>PROJ-123</code>)",
	)
}

func (uc *useCase) SendNotification(ctx context.Context, userID, message string) error {
	conn, err := uc.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil
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

func (uc *useCase) DeleteWebhook(ctx context.Context) error {
	return uc.tg.DeleteWebhook(ctx)
}

func (uc *useCase) BotInfo(ctx context.Context) *entity.TelegramBotInfo {
	return &entity.TelegramBotInfo{
		Username:   uc.botUsername,
		WebhookURL: uc.webhookURL,
	}
}

func (uc *useCase) StartPolling(ctx context.Context) {
	_ = uc.tg.Poll(ctx, func(ctx context.Context, update *entity.TelegramUpdate) {
		_ = uc.HandleUpdate(ctx, update)
	})
}
