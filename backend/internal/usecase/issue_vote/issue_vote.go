package issue_vote

import (
	"context"
	"fmt"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
)

type useCase struct {
	voteRepo  repository.IssueVoteRepository
	issueRepo repository.IssueRepository
}

func New(voteRepo repository.IssueVoteRepository, issueRepo repository.IssueRepository) UseCase {
	return &useCase{voteRepo: voteRepo, issueRepo: issueRepo}
}

func (uc *useCase) Toggle(ctx context.Context, issueID, actorID string) (bool, error) {
	if _, err := uc.issueRepo.GetByID(ctx, issueID); err != nil {
		return false, fmt.Errorf("issue_vote.Toggle: %w", err)
	}
	return uc.voteRepo.Toggle(ctx, issueID, actorID)
}

func (uc *useCase) GetSummary(ctx context.Context, issueID, actorID string) (*entity.IssueVoteSummary, error) {
	if _, err := uc.issueRepo.GetByID(ctx, issueID); err != nil {
		return nil, fmt.Errorf("issue_vote.GetSummary: %w", err)
	}
	return uc.voteRepo.GetSummary(ctx, issueID, actorID)
}
