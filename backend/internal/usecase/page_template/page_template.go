package page_template

import (
	"context"
	"fmt"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	repo repository.PageTemplateRepository
	log  logger.Logger
}

func New(repo repository.PageTemplateRepository, log logger.Logger) UseCase {
	return &useCase{repo: repo, log: log}
}

func (uc *useCase) Create(ctx context.Context, spaceID *string, createdBy string, req *entity.CreatePageTemplateReq) (*entity.PageTemplate, error) {
	content := req.Content
	if content == nil {
		content = map[string]any{}
	}
	t := &entity.PageTemplate{
		SpaceID:     spaceID,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Content:     content,
		ContentText: req.ContentText,
		Icon:        req.Icon,
		CreatedBy:   createdBy,
		IsGlobal:    req.IsGlobal,
	}
	if err := uc.repo.Create(ctx, t); err != nil {
		uc.log.Error(ctx, "pageTemplate.Create: db error", logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("pageTemplate.Create: %w", err)
	}
	uc.log.Info(ctx, "page template created", logger.String("id", t.ID))
	return t, nil
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.PageTemplate, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *useCase) List(ctx context.Context, filter *entity.PageTemplateFilter) ([]*entity.PageTemplate, int, error) {
	return uc.repo.List(ctx, filter)
}

func (uc *useCase) Update(ctx context.Context, id, actorID string, req *entity.UpdatePageTemplateReq) (*entity.PageTemplate, error) {
	t, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t.CreatedBy != actorID && !t.IsGlobal {
		return nil, apperr.Forbidden("only the creator can update this template")
	}

	if req.Name != nil {
		t.Name = *req.Name
	}
	if req.Description != nil {
		t.Description = req.Description
	}
	if req.Category != nil {
		t.Category = *req.Category
	}
	if req.Content != nil {
		t.Content = req.Content
	}
	if req.ContentText != nil {
		t.ContentText = *req.ContentText
	}
	if req.Icon != nil {
		t.Icon = req.Icon
	}

	if err := uc.repo.Update(ctx, t); err != nil {
		uc.log.Error(ctx, "pageTemplate.Update: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("pageTemplate.Update: %w", err)
	}
	return t, nil
}

func (uc *useCase) Delete(ctx context.Context, id, actorID string) error {
	t, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if t.CreatedBy != actorID {
		return apperr.Forbidden("only the creator can delete this template")
	}
	if err := uc.repo.Delete(ctx, id); err != nil {
		uc.log.Error(ctx, "pageTemplate.Delete: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return fmt.Errorf("pageTemplate.Delete: %w", err)
	}
	return nil
}
