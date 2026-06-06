package project_template

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
)

type UseCase interface {
	List(ctx context.Context) ([]*entity.ProjectTemplate, error)
	GetByID(ctx context.Context, id string) (*entity.ProjectTemplate, error)
}

type useCase struct {
	repo repository.ProjectTemplateRepository
}

func New(repo repository.ProjectTemplateRepository) UseCase {
	return &useCase{repo: repo}
}

func (uc *useCase) List(ctx context.Context) ([]*entity.ProjectTemplate, error) {
	return uc.repo.List(ctx)
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.ProjectTemplate, error) {
	return uc.repo.GetByID(ctx, id)
}
