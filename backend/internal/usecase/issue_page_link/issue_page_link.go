package issue_page_link

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
	repo      repository.IssuePageLinkRepository
	issueRepo repository.IssueRepository
	pageRepo  repository.PageRepository
	log       logger.Logger
}

func New(
	repo repository.IssuePageLinkRepository,
	issueRepo repository.IssueRepository,
	pageRepo repository.PageRepository,
	log logger.Logger,
) UseCase {
	return &useCase{repo: repo, issueRepo: issueRepo, pageRepo: pageRepo, log: log}
}

func (uc *useCase) Link(ctx context.Context, issueID, pageID, actorID string) (*entity.IssuePageLink, error) {
	if _, err := uc.issueRepo.GetByID(ctx, issueID); err != nil {
		return nil, apperr.NotFound("issue")
	}
	if _, err := uc.pageRepo.GetByID(ctx, pageID); err != nil {
		return nil, apperr.NotFound("page")
	}
	link := &entity.IssuePageLink{
		ID:        uuid.NewString(),
		IssueID:   issueID,
		PageID:    pageID,
		LinkedBy:  actorID,
		CreatedAt: time.Now().UTC(),
	}
	if err := uc.repo.Create(ctx, link); err != nil {
		uc.log.Error(ctx, "issue_page_link.Link: db error",
			logger.String("issue_id", issueID), logger.String("page_id", pageID),
			logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "issue-page linked", logger.String("issue_id", issueID), logger.String("page_id", pageID))
	return link, nil
}

func (uc *useCase) Unlink(ctx context.Context, issueID, pageID string) error {
	exists, err := uc.repo.Exists(ctx, issueID, pageID)
	if err != nil {
		return err
	}
	if !exists {
		return apperr.NotFound("issue-page link")
	}
	if err := uc.repo.Delete(ctx, issueID, pageID); err != nil {
		uc.log.Error(ctx, "issue_page_link.Unlink: db error",
			logger.String("issue_id", issueID), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "issue-page unlinked", logger.String("issue_id", issueID), logger.String("page_id", pageID))
	return nil
}

func (uc *useCase) ListByIssue(ctx context.Context, issueID string) ([]*entity.IssuePageLink, error) {
	return uc.repo.ListByIssue(ctx, issueID)
}

func (uc *useCase) ListByPage(ctx context.Context, pageID string) ([]*entity.IssuePageLink, error) {
	return uc.repo.ListByPage(ctx, pageID)
}
