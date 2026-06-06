package worker

import (
	"context"
	"fmt"
	"time"

	emailpkg "github.com/jira-backend/jiraflow-backend/internal/infrastructura/email"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type DailyDigestWorker struct {
	notifRepo   repository.NotificationRepository
	userRepo    repository.UserRepository
	emailSender emailpkg.Sender
	log         logger.Logger
}

func NewDailyDigestWorker(
	notifRepo repository.NotificationRepository,
	userRepo repository.UserRepository,
	emailSender emailpkg.Sender,
	log logger.Logger,
) *DailyDigestWorker {
	return &DailyDigestWorker{
		notifRepo:   notifRepo,
		userRepo:    userRepo,
		emailSender: emailSender,
		log:         log,
	}
}

// Run blocks until ctx is cancelled, firing the digest at atHourUTC every day.
func (w *DailyDigestWorker) Run(ctx context.Context, atHourUTC int) {
	for {
		next := nextOccurrence(atHourUTC)
		w.log.Info(ctx, fmt.Sprintf("daily digest: next run at %s", next.Format(time.RFC3339)))

		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Until(next)):
			w.send(ctx)
		}
	}
}

func (w *DailyDigestWorker) send(ctx context.Context) {
	isActive := true
	users, _, err := w.userRepo.List(ctx, &entity.UserFilter{
		Filter:   entity.Filter{Limit: 10000},
		IsActive: &isActive,
	})
	if err != nil {
		w.log.Error(ctx, "daily digest: list users", logger.SafeString("err", err.Error()))
		return
	}

	unread := true
	sent := 0
	for _, u := range users {
		if u.Email == "" {
			continue
		}
		notifs, _, err := w.notifRepo.ListByUser(ctx, u.ID, &entity.NotificationFilter{
			Filter: entity.Filter{Limit: 50},
			Unread: &unread,
		})
		if err != nil || len(notifs) == 0 {
			continue
		}

		subject := fmt.Sprintf("Your JiraFlow digest — %d unread notifications", len(notifs))
		body := buildDigestBody(u.FullName, notifs)

		if err := w.emailSender.SendRaw(ctx, []string{u.Email}, subject, body); err != nil {
			w.log.Error(ctx, "daily digest: send email",
				logger.String("user_id", u.ID),
				logger.SafeString("err", err.Error()),
			)
			continue
		}
		sent++
	}

	w.log.Info(ctx, fmt.Sprintf("daily digest: sent to %d users", sent))
}

func buildDigestBody(name string, notifs []*entity.Notification) string {
	body := fmt.Sprintf("Hi %s,\n\nHere's what happened in JiraFlow:\n\n", name)
	for _, n := range notifs {
		body += fmt.Sprintf("• [%s] %v\n", n.Type, n.Payload["title"])
	}
	body += "\nLog in to view details.\n\nYou're receiving this because you have email digests enabled.\n"
	return body
}

func nextOccurrence(hourUTC int) time.Time {
	now := time.Now().UTC()
	next := time.Date(now.Year(), now.Month(), now.Day(), hourUTC, 0, 0, 0, time.UTC)
	if !next.After(now) {
		next = next.Add(24 * time.Hour)
	}
	return next
}
