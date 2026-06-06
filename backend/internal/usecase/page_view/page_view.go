package page_view

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	repo repository.PageViewRepository
	log  logger.Logger
}

func New(repo repository.PageViewRepository, log logger.Logger) UseCase {
	return &useCase{repo: repo, log: log}
}

func (uc *useCase) RecordView(ctx context.Context, pageID string, userID *string, ipAddress *string) error {
	view := &entity.PageView{
		PageID:    pageID,
		UserID:    userID,
		IPAddress: ipAddress,
	}
	if err := uc.repo.Record(ctx, view); err != nil {
		uc.log.Error(ctx, "pageView.RecordView: db error", logger.String("page_id", pageID), logger.SafeString("err", err.Error()))
	}
	return nil
}

func (uc *useCase) GetAnalytics(ctx context.Context, pageID string) (*entity.PageAnalytics, error) {
	return uc.repo.GetAnalytics(ctx, pageID)
}

func (uc *useCase) ListRecentByUser(ctx context.Context, userID string, limit int) ([]*entity.RecentPage, error) {
	return uc.repo.ListRecentByUser(ctx, userID, limit)
}
