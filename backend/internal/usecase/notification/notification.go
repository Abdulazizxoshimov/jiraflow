package notification

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	notifRepo repository.NotificationRepository
	log       logger.Logger
}

func New(notifRepo repository.NotificationRepository, log logger.Logger) UseCase {
	return &useCase{notifRepo: notifRepo, log: log}
}

func (uc *useCase) Notify(ctx context.Context, n *entity.Notification) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	if n.CreatedAt.IsZero() {
		n.CreatedAt = time.Now().UTC()
	}
	if n.Payload == nil {
		n.Payload = map[string]any{}
	}
	if err := uc.notifRepo.Create(ctx, n); err != nil {
		uc.log.Error(ctx, "notification.Notify: db error", logger.String("user_id", n.UserID), logger.SafeString("err", err.Error()))
		return err
	}
	return nil
}

func (uc *useCase) ListByUser(ctx context.Context, userID string, filter *entity.NotificationFilter) ([]*entity.Notification, int, error) {
	return uc.notifRepo.ListByUser(ctx, userID, filter)
}

func (uc *useCase) MarkRead(ctx context.Context, userID string, ids []string) error {
	if err := uc.notifRepo.MarkRead(ctx, userID, ids); err != nil {
		uc.log.Error(ctx, "notification.MarkRead: db error", logger.String("user_id", userID), logger.SafeString("err", err.Error()))
		return err
	}
	return nil
}

func (uc *useCase) MarkAllRead(ctx context.Context, userID string) error {
	if err := uc.notifRepo.MarkAllRead(ctx, userID); err != nil {
		uc.log.Error(ctx, "notification.MarkAllRead: db error", logger.String("user_id", userID), logger.SafeString("err", err.Error()))
		return err
	}
	return nil
}

func (uc *useCase) Delete(ctx context.Context, id, userID string) error {
	if err := uc.notifRepo.Delete(ctx, id); err != nil {
		uc.log.Error(ctx, "notification.Delete: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return err
	}
	return nil
}

func (uc *useCase) CountUnread(ctx context.Context, userID string) (int, error) {
	return uc.notifRepo.CountUnread(ctx, userID)
}

func (uc *useCase) GetPreference(ctx context.Context, userID string) (*entity.NotificationPreference, error) {
	pref, err := uc.notifRepo.GetPreference(ctx, userID)
	if err != nil {
		return &entity.NotificationPreference{
			UserID:         userID,
			EmailAssigned:  true,
			EmailMentioned: true,
			EmailCommented: true,
			EmailStatus:    true,
			EmailWatcher:   false,
			DailyDigest:    false,
			UpdatedAt:      time.Now().UTC(),
		}, nil
	}
	return pref, nil
}

func (uc *useCase) UpdatePreference(ctx context.Context, userID string, req *entity.UpdateNotificationPreferenceReq) (*entity.NotificationPreference, error) {
	pref, _ := uc.notifRepo.GetPreference(ctx, userID)
	if pref == nil {
		pref = &entity.NotificationPreference{
			UserID:         userID,
			EmailAssigned:  true,
			EmailMentioned: true,
			EmailCommented: true,
			EmailStatus:    true,
		}
	}
	if req.EmailAssigned != nil {
		pref.EmailAssigned = *req.EmailAssigned
	}
	if req.EmailMentioned != nil {
		pref.EmailMentioned = *req.EmailMentioned
	}
	if req.EmailCommented != nil {
		pref.EmailCommented = *req.EmailCommented
	}
	if req.EmailStatus != nil {
		pref.EmailStatus = *req.EmailStatus
	}
	if req.EmailWatcher != nil {
		pref.EmailWatcher = *req.EmailWatcher
	}
	if req.DailyDigest != nil {
		pref.DailyDigest = *req.DailyDigest
	}
	if req.TelegramEnabled != nil {
		pref.TelegramEnabled = *req.TelegramEnabled
	}
	pref.UpdatedAt = time.Now().UTC()

	if err := uc.notifRepo.UpsertPreference(ctx, pref); err != nil {
		uc.log.Error(ctx, "notification.UpdatePreference: db error", logger.String("user_id", userID), logger.SafeString("err", err.Error()))
		return nil, err
	}
	return pref, nil
}
