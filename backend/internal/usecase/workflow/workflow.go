package workflow

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	repo repository.WorkflowRepository
	log  logger.Logger
}

func New(repo repository.WorkflowRepository, log logger.Logger) UseCase {
	return &useCase{repo: repo, log: log}
}

func (uc *useCase) Create(ctx context.Context, wf *entity.Workflow) (*entity.Workflow, error) {
	now := time.Now().UTC()
	wf.ID = uuid.NewString()
	wf.CreatedAt = now
	wf.UpdatedAt = now
	if err := uc.repo.Create(ctx, wf); err != nil {
		uc.log.Error(ctx, "workflow.Create: db error", logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "workflow created", logger.String("id", wf.ID))
	return wf, nil
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.Workflow, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *useCase) GetWithDetails(ctx context.Context, id string) (*entity.Workflow, error) {
	return uc.repo.GetWithDetails(ctx, id)
}

func (uc *useCase) List(ctx context.Context, filter *entity.Filter) ([]*entity.Workflow, int, error) {
	return uc.repo.List(ctx, filter)
}

func (uc *useCase) Update(ctx context.Context, id string, wf *entity.Workflow) (*entity.Workflow, error) {
	existing, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if wf.Name != "" {
		existing.Name = wf.Name
	}
	existing.Description = wf.Description
	if err := uc.repo.Update(ctx, existing); err != nil {
		uc.log.Error(ctx, "workflow.Update: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "workflow updated", logger.String("id", id))
	return existing, nil
}

func (uc *useCase) SetDefault(ctx context.Context, id string) error {
	if _, err := uc.repo.GetByID(ctx, id); err != nil {
		return err
	}
	if err := uc.repo.SetDefault(ctx, id); err != nil {
		uc.log.Error(ctx, "workflow.SetDefault: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "workflow set as default", logger.String("id", id))
	return nil
}

func (uc *useCase) Delete(ctx context.Context, id string) error {
	if err := uc.repo.SoftDelete(ctx, id); err != nil {
		uc.log.Error(ctx, "workflow.Delete: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "workflow deleted", logger.String("id", id))
	return nil
}

func (uc *useCase) CreateStatus(ctx context.Context, s *entity.WorkflowStatus) (*entity.WorkflowStatus, error) {
	now := time.Now().UTC()
	s.ID = uuid.NewString()
	s.CreatedAt = now
	s.UpdatedAt = now
	if err := uc.repo.CreateStatus(ctx, s); err != nil {
		uc.log.Error(ctx, "workflow.CreateStatus: db error", logger.String("workflow_id", s.WorkflowID), logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "workflow status created", logger.String("id", s.ID))
	return s, nil
}

func (uc *useCase) UpdateStatus(ctx context.Context, s *entity.WorkflowStatus) (*entity.WorkflowStatus, error) {
	existing, err := uc.repo.GetStatusByID(ctx, s.ID)
	if err != nil {
		return nil, err
	}
	if s.Name != "" {
		existing.Name = s.Name
	}
	if s.Category != "" {
		existing.Category = s.Category
	}
	if s.Color != "" {
		existing.Color = s.Color
	}
	existing.Position = s.Position
	if err := uc.repo.UpdateStatus(ctx, existing); err != nil {
		uc.log.Error(ctx, "workflow.UpdateStatus: db error", logger.String("id", s.ID), logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "workflow status updated", logger.String("id", s.ID))
	return existing, nil
}

func (uc *useCase) DeleteStatus(ctx context.Context, id string) error {
	if err := uc.repo.DeleteStatus(ctx, id); err != nil {
		uc.log.Error(ctx, "workflow.DeleteStatus: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "workflow status deleted", logger.String("id", id))
	return nil
}

func (uc *useCase) CreateTransition(ctx context.Context, t *entity.WorkflowTransition) (*entity.WorkflowTransition, error) {
	t.ID = uuid.NewString()
	t.CreatedAt = time.Now().UTC()
	if err := uc.repo.CreateTransition(ctx, t); err != nil {
		uc.log.Error(ctx, "workflow.CreateTransition: db error", logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "workflow transition created", logger.String("id", t.ID))
	return t, nil
}

func (uc *useCase) DeleteTransition(ctx context.Context, id string) error {
	if err := uc.repo.DeleteTransition(ctx, id); err != nil {
		uc.log.Error(ctx, "workflow.DeleteTransition: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "workflow transition deleted", logger.String("id", id))
	return nil
}
