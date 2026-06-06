package label

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	labelRepo repository.LabelRepository
	log       logger.Logger
}

func New(labelRepo repository.LabelRepository, log logger.Logger) UseCase {
	return &useCase{labelRepo: labelRepo, log: log}
}

func (uc *useCase) Create(ctx context.Context, projectID string, req *entity.CreateLabelReq) (*entity.Label, error) {
	color := req.Color
	if color == "" {
		color = "#6B7280"
	}
	l := &entity.Label{
		ID:        uuid.NewString(),
		ProjectID: projectID,
		Name:      req.Name,
		Color:     color,
		CreatedAt: time.Now().UTC(),
	}
	if err := uc.labelRepo.Create(ctx, l); err != nil {
		uc.log.Error(ctx, "label.Create: db error", logger.String("project_id", projectID), logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "label created", logger.String("id", l.ID), logger.String("project_id", projectID))
	return l, nil
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.Label, error) {
	return uc.labelRepo.GetByID(ctx, id)
}

func (uc *useCase) ListByProject(ctx context.Context, projectID string) ([]*entity.Label, error) {
	return uc.labelRepo.ListByProject(ctx, projectID)
}

func (uc *useCase) Update(ctx context.Context, id string, req *entity.UpdateLabelReq) (*entity.Label, error) {
	l, err := uc.labelRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		l.Name = *req.Name
	}
	if req.Color != nil {
		l.Color = *req.Color
	}
	if err := uc.labelRepo.Update(ctx, l); err != nil {
		uc.log.Error(ctx, "label.Update: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "label updated", logger.String("id", id))
	return l, nil
}

func (uc *useCase) Delete(ctx context.Context, id string) error {
	if err := uc.labelRepo.Delete(ctx, id); err != nil {
		uc.log.Error(ctx, "label.Delete: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "label deleted", logger.String("id", id))
	return nil
}
