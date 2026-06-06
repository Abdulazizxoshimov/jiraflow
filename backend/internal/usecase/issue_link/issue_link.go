package issue_link

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	linkRepo  repository.IssueLinkRepository
	issueRepo repository.IssueRepository
	log       logger.Logger
}

func New(linkRepo repository.IssueLinkRepository, issueRepo repository.IssueRepository, log logger.Logger) UseCase {
	return &useCase{linkRepo: linkRepo, issueRepo: issueRepo, log: log}
}

func (uc *useCase) Create(ctx context.Context, link *entity.IssueLink) (*entity.IssueLink, error) {
	if _, err := uc.issueRepo.GetByID(ctx, link.SourceID); err != nil {
		return nil, apperr.NotFound("source issue")
	}
	if _, err := uc.issueRepo.GetByID(ctx, link.TargetID); err != nil {
		return nil, apperr.NotFound("target issue")
	}
	if link.ID == "" {
		link.ID = uuid.NewString()
	}
	if link.CreatedAt.IsZero() {
		link.CreatedAt = time.Now().UTC()
	}
	if err := uc.linkRepo.Create(ctx, link); err != nil {
		uc.log.Error(ctx, "issue_link.Create: db error", logger.String("source_id", link.SourceID), logger.String("target_id", link.TargetID), logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "issue link created", logger.String("id", link.ID))
	return link, nil
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.IssueLink, error) {
	return uc.linkRepo.GetByID(ctx, id)
}

func (uc *useCase) ListByIssue(ctx context.Context, issueID string) ([]*entity.IssueLink, error) {
	return uc.linkRepo.ListByIssue(ctx, issueID)
}

func (uc *useCase) Delete(ctx context.Context, id, actorID string) error {
	if err := uc.linkRepo.Delete(ctx, id); err != nil {
		uc.log.Error(ctx, "issue_link.Delete: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "issue link deleted", logger.String("id", id))
	return nil
}
