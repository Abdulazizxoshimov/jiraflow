package page_macro

import (
	"context"
	"fmt"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
)

type useCase struct {
	repo repository.PageMacroRepository
}

func New(repo repository.PageMacroRepository) UseCase {
	return &useCase{repo: repo}
}

func (uc *useCase) Upsert(ctx context.Context, pageID string, req *entity.UpsertPageMacroReq) (*entity.PageMacro, error) {
	macro := &entity.PageMacro{
		PageID:    pageID,
		MacroType: req.MacroType,
		Config:    req.Config,
	}
	if err := uc.repo.Upsert(ctx, macro); err != nil {
		return nil, fmt.Errorf("pageMacro.Upsert: %w", err)
	}
	return macro, nil
}

func (uc *useCase) ListByPage(ctx context.Context, pageID string) ([]*entity.PageMacro, error) {
	return uc.repo.ListByPage(ctx, pageID)
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.PageMacro, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *useCase) Delete(ctx context.Context, id string) error {
	if _, err := uc.repo.GetByID(ctx, id); err != nil {
		return err
	}
	return uc.repo.Delete(ctx, id)
}
