package activity_feed

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	repo repository.ActivityFeedRepository
	log  logger.Logger
}

func New(repo repository.ActivityFeedRepository, log logger.Logger) UseCase {
	return &useCase{repo: repo, log: log}
}

func (uc *useCase) Record(ctx context.Context, ev *entity.ActivityEvent) error {
	if ev.ID == "" {
		ev.ID = uuid.NewString()
	}
	if ev.CreatedAt.IsZero() {
		ev.CreatedAt = time.Now().UTC()
	}
	if err := uc.repo.Create(ctx, ev); err != nil {
		uc.log.Error(ctx, "activity_feed.Record: db error",
			logger.String("action", ev.Action),
			logger.String("entity_type", ev.EntityType),
			logger.SafeString("err", err.Error()))
		return err
	}
	return nil
}

func (uc *useCase) List(ctx context.Context, filter *entity.ActivityFilter) ([]*entity.ActivityEvent, int, error) {
	return uc.repo.List(ctx, filter)
}
