package content_property

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
)

type useCase struct {
	repo repository.ContentPropertyRepository
}

func New(repo repository.ContentPropertyRepository) UseCase {
	return &useCase{repo: repo}
}

func (uc *useCase) Set(ctx context.Context, entityType, entityID, key string, value map[string]any) (*entity.ContentProperty, error) {
	p := &entity.ContentProperty{
		EntityType: entityType,
		EntityID:   entityID,
		Key:        key,
		Value:      value,
	}
	if err := uc.repo.Set(ctx, p); err != nil {
		return nil, err
	}
	return uc.repo.Get(ctx, entityType, entityID, key)
}

func (uc *useCase) Get(ctx context.Context, entityType, entityID, key string) (*entity.ContentProperty, error) {
	return uc.repo.Get(ctx, entityType, entityID, key)
}

func (uc *useCase) List(ctx context.Context, entityType, entityID string) ([]*entity.ContentProperty, error) {
	return uc.repo.List(ctx, entityType, entityID)
}

func (uc *useCase) Delete(ctx context.Context, entityType, entityID, key string) error {
	return uc.repo.Delete(ctx, entityType, entityID, key)
}
