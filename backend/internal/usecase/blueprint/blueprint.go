package blueprint

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
)

type useCase struct {
	repo     repository.BlueprintRepository
	pageRepo repository.PageRepository
}

func New(repo repository.BlueprintRepository, pageRepo repository.PageRepository) UseCase {
	return &useCase{repo: repo, pageRepo: pageRepo}
}

func (uc *useCase) Create(ctx context.Context, req *entity.CreateBlueprintReq) (*entity.Blueprint, error) {
	b := &entity.Blueprint{
		ID:           uuid.NewString(),
		Name:         req.Name,
		Description:  req.Description,
		IconURL:      req.IconURL,
		Category:     req.Category,
		TemplateBody: req.TemplateBody,
		Schema:       req.Schema,
		IsSystem:     false,
	}
	if err := uc.repo.Create(ctx, b); err != nil {
		return nil, fmt.Errorf("blueprint.Create: %w", err)
	}
	return b, nil
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.Blueprint, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *useCase) List(ctx context.Context) ([]*entity.Blueprint, error) {
	return uc.repo.List(ctx)
}

func (uc *useCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *useCase) CreatePage(ctx context.Context, blueprintID, actorID string, req *entity.CreatePageFromBlueprintReq) (*entity.Page, error) {
	b, err := uc.repo.GetByID(ctx, blueprintID)
	if err != nil {
		return nil, err
	}

	contentText := ""
	if b.TemplateBody != nil {
		contentText = *b.TemplateBody
	}

	now := time.Now().UTC()
	p := &entity.Page{
		ID:             uuid.NewString(),
		SpaceID:        req.SpaceID,
		ParentID:       req.ParentID,
		Title:          req.Title,
		Content:        map[string]any{"type": "doc", "content": []any{}},
		ContentText:    contentText,
		AuthorID:       actorID,
		LastEditorID:   actorID,
		CurrentVersion: 1,
		Status:         "draft",
		Position:       1,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if err := uc.pageRepo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("blueprint.CreatePage: %w", err)
	}
	return p, nil
}
