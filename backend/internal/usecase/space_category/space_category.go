package space_category

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
)

type useCase struct {
	repo repository.SpaceCategoryRepository
}

func New(repo repository.SpaceCategoryRepository) UseCase {
	return &useCase{repo: repo}
}

func (uc *useCase) Create(ctx context.Context, req *entity.CreateSpaceCategoryReq) (*entity.SpaceCategory, error) {
	c := &entity.SpaceCategory{
		ID:    uuid.NewString(),
		Name:  req.Name,
		Color: req.Color,
	}
	if err := uc.repo.Create(ctx, c); err != nil {
		return nil, fmt.Errorf("space_category.Create: %w", err)
	}
	return c, nil
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.SpaceCategory, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *useCase) List(ctx context.Context) ([]*entity.SpaceCategory, error) {
	return uc.repo.List(ctx)
}

func (uc *useCase) Update(ctx context.Context, id string, req *entity.UpdateSpaceCategoryReq) (*entity.SpaceCategory, error) {
	c, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		c.Name = req.Name
	}
	c.Color = req.Color
	if err := uc.repo.Update(ctx, c); err != nil {
		return nil, fmt.Errorf("space_category.Update: %w", err)
	}
	return c, nil
}

func (uc *useCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
