package worker

import (
	"context"
	"time"

	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

// NewTokenCleanupWorker deletes expired refresh tokens, password resets,
// and invites from the database every hour.
func NewTokenCleanupWorker(
	authRepo repository.AuthRepository,
	inviteRepo repository.InviteRepository,
	log logger.Logger,
) *BackgroundWorker {
	jobs := []Job{
		{
			Name:     "token-cleanup",
			Interval: time.Hour,
			Run: func(ctx context.Context) {
				if err := authRepo.DeleteExpiredTokens(ctx); err != nil {
					log.Error(ctx, "token cleanup failed", logger.Error(err))
				}
				if err := authRepo.DeleteExpiredPasswordResets(ctx); err != nil {
					log.Error(ctx, "password reset cleanup failed", logger.Error(err))
				}
				if err := inviteRepo.DeleteExpired(ctx); err != nil {
					log.Error(ctx, "invite cleanup failed", logger.Error(err))
				}
			},
		},
	}
	return NewBackgroundWorker("token-cleanup-worker", jobs...)
}
