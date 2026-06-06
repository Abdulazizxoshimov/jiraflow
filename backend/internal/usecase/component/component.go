package component

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	repo repository.ComponentRepository
	log  logger.Logger
}

func New(repo repository.ComponentRepository, log logger.Logger) UseCase {
	return &useCase{repo: repo, log: log}
}

func (uc *useCase) Create(ctx context.Context, projectID string, req *entity.CreateComponentReq) (*entity.Component, error) {
	now := time.Now().UTC()
	c := &entity.Component{
		ID:          uuid.NewString(),
		ProjectID:   projectID,
		Name:        req.Name,
		Description: req.Description,
		LeadID:      req.LeadID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := uc.repo.Create(ctx, c); err != nil {
		uc.log.Error(ctx, "component.Create: db error", logger.String("project_id", projectID), logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "component created", logger.String("id", c.ID))
	return uc.repo.GetByID(ctx, c.ID)
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.Component, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *useCase) List(ctx context.Context, projectID string) ([]*entity.Component, error) {
	return uc.repo.List(ctx, projectID)
}

func (uc *useCase) Update(ctx context.Context, id string, req *entity.UpdateComponentReq) (*entity.Component, error) {
	c, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		c.Name = *req.Name
	}
	if req.Description != nil {
		c.Description = req.Description
	}
	if req.LeadID != nil {
		c.LeadID = req.LeadID
	}
	if err := uc.repo.Update(ctx, c); err != nil {
		uc.log.Error(ctx, "component.Update: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return nil, err
	}
	return uc.repo.GetByID(ctx, id)
}

func (uc *useCase) Delete(ctx context.Context, id string) error {
	if _, err := uc.repo.GetByID(ctx, id); err != nil {
		return err
	}
	return uc.repo.Delete(ctx, id)
}
