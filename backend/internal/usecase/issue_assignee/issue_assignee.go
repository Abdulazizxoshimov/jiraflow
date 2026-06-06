package issue_assignee

import (
	"context"
	"fmt"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	repo      repository.IssueAssigneeRepository
	issueRepo repository.IssueRepository
	log       logger.Logger
}

func New(repo repository.IssueAssigneeRepository, issueRepo repository.IssueRepository, log logger.Logger) UseCase {
	return &useCase{repo: repo, issueRepo: issueRepo, log: log}
}

func (uc *useCase) Set(ctx context.Context, issueID string, req *entity.SetIssueAssigneesReq) error {
	if _, err := uc.issueRepo.GetByID(ctx, issueID); err != nil {
		return err
	}
	if err := uc.repo.Set(ctx, issueID, req.UserIDs, req.PrimaryID); err != nil {
		uc.log.Error(ctx, "issueAssignee.Set: db error", logger.String("issue_id", issueID), logger.SafeString("err", err.Error()))
		return fmt.Errorf("issueAssignee.Set: %w", err)
	}
	return nil
}

func (uc *useCase) List(ctx context.Context, issueID string) ([]*entity.IssueAssignee, error) {
	return uc.repo.List(ctx, issueID)
}

func (uc *useCase) Remove(ctx context.Context, issueID, userID string) error {
	if err := uc.repo.Remove(ctx, issueID, userID); err != nil {
		uc.log.Error(ctx, "issueAssignee.Remove: db error", logger.String("issue_id", issueID), logger.SafeString("err", err.Error()))
		return fmt.Errorf("issueAssignee.Remove: %w", err)
	}
	return nil
}
