package worklog

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	worklogRepo repository.WorklogRepository
	issueRepo   repository.IssueRepository
	log         logger.Logger
}

func New(worklogRepo repository.WorklogRepository, issueRepo repository.IssueRepository, log logger.Logger) UseCase {
	return &useCase{worklogRepo: worklogRepo, issueRepo: issueRepo, log: log}
}

func (uc *useCase) Create(ctx context.Context, issueID, userID string, req *entity.CreateWorklogReq) (*entity.Worklog, error) {
	// Issue mavjudligini tekshirish
	if _, err := uc.issueRepo.GetByID(ctx, issueID); err != nil {
		return nil, err
	}

	startedAt := req.StartedAt
	if startedAt.IsZero() {
		startedAt = time.Now().UTC()
	}

	now := time.Now().UTC()
	w := &entity.Worklog{
		ID:          uuid.NewString(),
		IssueID:     issueID,
		UserID:      userID,
		TimeSpent:   req.TimeSpent,
		StartedAt:   startedAt,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := uc.worklogRepo.Create(ctx, w); err != nil {
		uc.log.Error(ctx, "worklog.Create: db error", logger.String("issue_id", issueID), logger.SafeString("err", err.Error()))
		return nil, err
	}

	// remaining_estimate'ni avtomatik kamaytirish
	go func() {
		issue, err := uc.issueRepo.GetByID(ctx, issueID)
		if err != nil || issue.RemainingEstimate == nil {
			return
		}
		newRemaining := *issue.RemainingEstimate - req.TimeSpent
		if newRemaining < 0 {
			newRemaining = 0
		}
		_ = uc.issueRepo.UpdateEstimates(ctx, issueID, nil, &newRemaining)
	}()

	uc.log.Info(ctx, "worklog created", logger.String("id", w.ID), logger.String("issue_id", issueID))
	return uc.worklogRepo.GetByID(ctx, w.ID)
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.Worklog, error) {
	return uc.worklogRepo.GetByID(ctx, id)
}

func (uc *useCase) List(ctx context.Context, filter *entity.WorklogFilter) ([]*entity.Worklog, int, error) {
	return uc.worklogRepo.List(ctx, filter)
}

func (uc *useCase) Update(ctx context.Context, id, actorID string, req *entity.UpdateWorklogReq) (*entity.Worklog, error) {
	w, err := uc.worklogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if w.UserID != actorID {
		return nil, apperr.Forbidden("only the author can edit this worklog")
	}

	if req.TimeSpent != nil {
		w.TimeSpent = *req.TimeSpent
	}
	if req.StartedAt != nil {
		w.StartedAt = *req.StartedAt
	}
	if req.Description != nil {
		w.Description = req.Description
	}

	if err := uc.worklogRepo.Update(ctx, w); err != nil {
		uc.log.Error(ctx, "worklog.Update: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("worklog update: %w", err)
	}
	return uc.worklogRepo.GetByID(ctx, id)
}

func (uc *useCase) Delete(ctx context.Context, id, actorID string) error {
	w, err := uc.worklogRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if w.UserID != actorID {
		return apperr.Forbidden("only the author can delete this worklog")
	}
	return uc.worklogRepo.Delete(ctx, id)
}

func (uc *useCase) GetTimeSummary(ctx context.Context, issueID string) (*entity.TimeSpentSummary, error) {
	issue, err := uc.issueRepo.GetByID(ctx, issueID)
	if err != nil {
		return nil, err
	}
	total, err := uc.worklogRepo.SumByIssue(ctx, issueID)
	if err != nil {
		return nil, err
	}
	return &entity.TimeSpentSummary{
		IssueID:           issueID,
		OriginalEstimate:  issue.OriginalEstimate,
		RemainingEstimate: issue.RemainingEstimate,
		TimeSpentTotal:    total,
	}, nil
}

func (uc *useCase) UpdateEstimates(ctx context.Context, issueID string, original, remaining *int) error {
	if err := uc.issueRepo.UpdateEstimates(ctx, issueID, original, remaining); err != nil {
		uc.log.Error(ctx, "worklog.UpdateEstimates: db error", logger.String("issue_id", issueID), logger.SafeString("err", err.Error()))
		return fmt.Errorf("update estimates: %w", err)
	}
	return nil
}
