package version

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	repo repository.VersionRepository
	log  logger.Logger
}

func New(repo repository.VersionRepository, log logger.Logger) UseCase {
	return &useCase{repo: repo, log: log}
}

func (uc *useCase) Create(ctx context.Context, projectID string, req *entity.CreateVersionReq) (*entity.Version, error) {
	now := time.Now().UTC()
	v := &entity.Version{
		ID:          uuid.NewString(),
		ProjectID:   projectID,
		Name:        req.Name,
		Description: req.Description,
		Status:      "unreleased",
		StartDate:   req.StartDate,
		ReleaseDate: req.ReleaseDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := uc.repo.Create(ctx, v); err != nil {
		uc.log.Error(ctx, "version.Create: db error", logger.String("project_id", projectID), logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "version created", logger.String("id", v.ID))
	return uc.repo.GetByID(ctx, v.ID)
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.Version, error) {
	v, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	total, done, err := uc.repo.GetProgress(ctx, id)
	if err == nil {
		v.IssueCount = total
		v.DoneCount = done
		if total > 0 {
			v.Progress = done * 100 / total
		}
	}
	return v, nil
}

func (uc *useCase) List(ctx context.Context, projectID string) ([]*entity.Version, error) {
	versions, err := uc.repo.List(ctx, projectID)
	if err != nil {
		return nil, err
	}
	for _, v := range versions {
		total, done, err := uc.repo.GetProgress(ctx, v.ID)
		if err == nil {
			v.IssueCount = total
			v.DoneCount = done
			if total > 0 {
				v.Progress = done * 100 / total
			}
		}
	}
	return versions, nil
}

func (uc *useCase) Update(ctx context.Context, id string, req *entity.UpdateVersionReq) (*entity.Version, error) {
	v, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		v.Name = *req.Name
	}
	if req.Description != nil {
		v.Description = req.Description
	}
	if req.StartDate != nil {
		v.StartDate = req.StartDate
	}
	if req.ReleaseDate != nil {
		v.ReleaseDate = req.ReleaseDate
	}
	if err := uc.repo.Update(ctx, v); err != nil {
		uc.log.Error(ctx, "version.Update: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return nil, err
	}
	return uc.GetByID(ctx, id)
}

func (uc *useCase) Release(ctx context.Context, id string, req *entity.ReleaseVersionReq) (*entity.Version, error) {
	releasedAt := time.Now().UTC()
	if req.ReleasedAt != nil {
		releasedAt = *req.ReleasedAt
	}
	if err := uc.repo.Release(ctx, id, releasedAt); err != nil {
		uc.log.Error(ctx, "version.Release: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "version released", logger.String("id", id))
	return uc.GetByID(ctx, id)
}

func (uc *useCase) Archive(ctx context.Context, id string) error {
	return uc.repo.Archive(ctx, id)
}

func (uc *useCase) Delete(ctx context.Context, id string) error {
	if _, err := uc.repo.GetByID(ctx, id); err != nil {
		return err
	}
	return uc.repo.Delete(ctx, id)
}
