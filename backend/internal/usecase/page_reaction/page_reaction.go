package page_reaction

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	repo repository.PageReactionRepository
	log  logger.Logger
}

func New(repo repository.PageReactionRepository, log logger.Logger) UseCase {
	return &useCase{repo: repo, log: log}
}

func (uc *useCase) Toggle(ctx context.Context, pageID, userID, emoji string) (bool, error) {
	added, err := uc.repo.Toggle(ctx, pageID, userID, emoji)
	if err != nil {
		uc.log.Error(ctx, "pageReaction.Toggle: db error",
			logger.String("page_id", pageID),
			logger.SafeString("err", err.Error()),
		)
		return false, err
	}
	return added, nil
}

func (uc *useCase) ListByPage(ctx context.Context, pageID, viewerUserID string) ([]*entity.PageReactionSummary, error) {
	return uc.repo.ListByPage(ctx, pageID, viewerUserID)
}

func (uc *useCase) ListUsers(ctx context.Context, pageID, emoji string) ([]*entity.PageReaction, error) {
	return uc.repo.ListUsers(ctx, pageID, emoji)
}
