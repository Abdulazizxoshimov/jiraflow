package project

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	repo      repository.ProjectRepository
	workflow  repository.WorkflowRepository
	spaceRepo repository.SpaceRepository
	log       logger.Logger
}

func New(repo repository.ProjectRepository, workflow repository.WorkflowRepository, spaceRepo repository.SpaceRepository, log logger.Logger) UseCase {
	return &useCase{repo: repo, workflow: workflow, spaceRepo: spaceRepo, log: log}
}

func (uc *useCase) Create(ctx context.Context, p *entity.Project, actorID string) (*entity.Project, error) {
	wf, err := uc.workflow.GetDefault(ctx)
	if err != nil {
		uc.log.Error(ctx, "project.Create: get default workflow failed", logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("project.Create get default workflow: %w", err)
	}
	now := time.Now().UTC()
	p.ID = uuid.NewString()
	p.WorkflowID = wf.ID
	p.LeadID = actorID
	p.CreatedAt = now
	p.UpdatedAt = now
	if err := uc.repo.Create(ctx, p); err != nil {
		uc.log.Error(ctx, "project.Create: db error", logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "project created", logger.String("id", p.ID), logger.String("lead_id", actorID))

	// Project uchun avtomatik Confluence space yaratish
	go uc.autoCreateSpace(context.Background(), p, actorID)

	return p, nil
}

func (uc *useCase) autoCreateSpace(ctx context.Context, p *entity.Project, leadID string) {
	spaceKey := strings.ToUpper(p.Key)
	now := time.Now().UTC()
	s := &entity.Space{
		ID:        uuid.NewString(),
		Key:       spaceKey,
		Name:      p.Name,
		Type:      "project",
		LeadID:    leadID,
		ProjectID: &p.ID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := uc.spaceRepo.Create(ctx, s); err != nil {
		uc.log.Warn(ctx, "project.autoCreateSpace: failed",
			logger.String("project_id", p.ID), logger.SafeString("err", err.Error()))
	}
}

func (uc *useCase) GetLinkedSpace(ctx context.Context, projectID string) (*entity.Space, error) {
	return uc.spaceRepo.GetByProjectID(ctx, projectID)
}

func (uc *useCase) GetDashboard(ctx context.Context, projectID string) (*entity.ProjectDashboard, error) {
	if _, err := uc.repo.GetByID(ctx, projectID); err != nil {
		return nil, err
	}
	d, err := uc.repo.GetDashboard(ctx, projectID)
	if err != nil {
		uc.log.Error(ctx, "project.GetDashboard: db error",
			logger.String("project_id", projectID), logger.SafeString("err", err.Error()))
		return nil, err
	}
	return d, nil
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.Project, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *useCase) GetByKey(ctx context.Context, key string) (*entity.Project, error) {
	return uc.repo.GetByKey(ctx, key)
}

func (uc *useCase) List(ctx context.Context, filter *entity.ProjectFilter) ([]*entity.Project, int, error) {
	return uc.repo.List(ctx, filter)
}

func (uc *useCase) Update(ctx context.Context, id string, p *entity.Project, actorID string) (*entity.Project, error) {
	existing, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	existing.Name = p.Name
	existing.Description = p.Description
	existing.IconURL = p.IconURL
	if p.LeadID != "" {
		existing.LeadID = p.LeadID
	}
	if p.WorkflowID != "" {
		existing.WorkflowID = p.WorkflowID
	}
	if err := uc.repo.Update(ctx, existing); err != nil {
		uc.log.Error(ctx, "project.Update: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "project updated", logger.String("id", id))
	return existing, nil
}

func (uc *useCase) Archive(ctx context.Context, id string, actorID string) error {
	p, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if p.LeadID != actorID {
		uc.log.Warn(ctx, "project.Archive: forbidden", logger.String("id", id), logger.String("actor_id", actorID))
		return apperr.Forbidden("only the project lead can archive this project")
	}
	p.IsArchived = true
	if err := uc.repo.Update(ctx, p); err != nil {
		uc.log.Error(ctx, "project.Archive: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "project archived", logger.String("id", id))
	return nil
}

func (uc *useCase) Delete(ctx context.Context, id string, actorID string) error {
	p, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if p.LeadID != actorID {
		uc.log.Warn(ctx, "project.Delete: forbidden", logger.String("id", id), logger.String("actor_id", actorID))
		return apperr.Forbidden("only the project lead can delete this project")
	}
	if err := uc.repo.SoftDelete(ctx, id); err != nil {
		uc.log.Error(ctx, "project.Delete: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "project deleted", logger.String("id", id))
	return nil
}
