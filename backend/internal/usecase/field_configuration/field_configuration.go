package field_configuration

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
)

type UseCase interface {
	Create(ctx context.Context, req *entity.CreateFieldConfigurationReq) (*entity.FieldConfiguration, error)
	GetByID(ctx context.Context, id string) (*entity.FieldConfiguration, error)
	List(ctx context.Context, projectID *string) ([]*entity.FieldConfiguration, error)
	Delete(ctx context.Context, id string) error
}

type useCase struct {
	repo repository.FieldConfigurationRepository
}

func New(repo repository.FieldConfigurationRepository) UseCase {
	return &useCase{repo: repo}
}

func (uc *useCase) Create(ctx context.Context, req *entity.CreateFieldConfigurationReq) (*entity.FieldConfiguration, error) {
	c := &entity.FieldConfiguration{
		ID:        uuid.NewString(),
		Name:      req.Name,
		ProjectID: req.ProjectID,
	}
	for _, item := range req.Items {
		c.Items = append(c.Items, &entity.FieldConfigItem{
			ID:          uuid.NewString(),
			ConfigID:    c.ID,
			FieldName:   item.FieldName,
			IsRequired:  item.IsRequired,
			IsHidden:    item.IsHidden,
			Description: item.Description,
		})
	}
	if err := uc.repo.Create(ctx, c); err != nil {
		return nil, fmt.Errorf("field_configuration.Create: %w", err)
	}
	return uc.repo.GetByID(ctx, c.ID)
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.FieldConfiguration, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *useCase) List(ctx context.Context, projectID *string) ([]*entity.FieldConfiguration, error) {
	return uc.repo.List(ctx, projectID)
}

func (uc *useCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
