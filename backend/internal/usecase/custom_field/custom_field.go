package custom_field

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	cfRepo repository.CustomFieldRepository
	log    logger.Logger
}

func New(cfRepo repository.CustomFieldRepository, log logger.Logger) UseCase {
	return &useCase{cfRepo: cfRepo, log: log}
}

func (uc *useCase) Create(ctx context.Context, projectID string, cf *entity.CustomField) (*entity.CustomField, error) {
	existing, _ := uc.cfRepo.GetByKey(ctx, projectID, cf.FieldKey)
	if existing != nil {
		return nil, apperr.Conflict("field key already exists in this project")
	}

	now := time.Now().UTC()
	cf.ID = uuid.NewString()
	cf.ProjectID = projectID
	cf.CreatedAt = now
	cf.UpdatedAt = now

	if err := uc.cfRepo.Create(ctx, cf); err != nil {
		uc.log.Error(ctx, "custom_field.Create: db error", logger.String("project_id", projectID), logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "custom field created", logger.String("id", cf.ID), logger.String("project_id", projectID))
	return cf, nil
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.CustomField, error) {
	return uc.cfRepo.GetByID(ctx, id)
}

func (uc *useCase) ListByProject(ctx context.Context, projectID string) ([]*entity.CustomField, error) {
	return uc.cfRepo.ListByProject(ctx, projectID)
}

func (uc *useCase) Update(ctx context.Context, id string, cf *entity.CustomField) (*entity.CustomField, error) {
	existing, err := uc.cfRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cf.Name != "" {
		existing.Name = cf.Name
	}
	if cf.Options != nil {
		existing.Options = cf.Options
	}
	existing.IsRequired = cf.IsRequired
	existing.UpdatedAt = time.Now().UTC()

	if err := uc.cfRepo.Update(ctx, existing); err != nil {
		uc.log.Error(ctx, "custom_field.Update: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "custom field updated", logger.String("id", id))
	return existing, nil
}

func (uc *useCase) Delete(ctx context.Context, id string) error {
	if err := uc.cfRepo.Delete(ctx, id); err != nil {
		uc.log.Error(ctx, "custom_field.Delete: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "custom field deleted", logger.String("id", id))
	return nil
}

func (uc *useCase) Reorder(ctx context.Context, projectID string, positions map[string]int) error {
	return uc.cfRepo.ReorderFields(ctx, projectID, positions)
}
