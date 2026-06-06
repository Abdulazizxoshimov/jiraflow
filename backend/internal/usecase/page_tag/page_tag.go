package page_tag

import (
	"context"
	"fmt"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	repo repository.PageTagRepository
	log  logger.Logger
}

func New(repo repository.PageTagRepository, log logger.Logger) UseCase {
	return &useCase{repo: repo, log: log}
}

func (uc *useCase) Create(ctx context.Context, spaceID string, req *entity.CreatePageTagReq) (*entity.PageTag, error) {
	tag := &entity.PageTag{
		SpaceID: spaceID,
		Name:    req.Name,
		Color:   req.Color,
	}
	if tag.Color == "" {
		tag.Color = "#6B7280"
	}
	if err := uc.repo.Create(ctx, tag); err != nil {
		uc.log.Error(ctx, "pageTag.Create: db error", logger.String("space_id", spaceID), logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("pageTag.Create: %w", err)
	}
	uc.log.Info(ctx, "page tag created", logger.String("id", tag.ID))
	return tag, nil
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.PageTag, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *useCase) List(ctx context.Context, spaceID string) ([]*entity.PageTag, error) {
	return uc.repo.List(ctx, spaceID)
}

func (uc *useCase) Update(ctx context.Context, id string, req *entity.UpdatePageTagReq) (*entity.PageTag, error) {
	tag, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		tag.Name = *req.Name
	}
	if req.Color != nil {
		tag.Color = *req.Color
	}
	if err := uc.repo.Update(ctx, tag); err != nil {
		uc.log.Error(ctx, "pageTag.Update: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("pageTag.Update: %w", err)
	}
	return tag, nil
}

func (uc *useCase) Delete(ctx context.Context, id string) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		uc.log.Error(ctx, "pageTag.Delete: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return fmt.Errorf("pageTag.Delete: %w", err)
	}
	return nil
}

func (uc *useCase) SetPageTags(ctx context.Context, pageID string, req *entity.SetPageTagsReq) error {
	return uc.repo.SetPageTags(ctx, pageID, req.TagIDs)
}

func (uc *useCase) GetPageTags(ctx context.Context, pageID string) ([]*entity.PageTag, error) {
	return uc.repo.GetPageTags(ctx, pageID)
}

func (uc *useCase) GetPagesByTag(ctx context.Context, tagID string, filter *entity.Filter) ([]*entity.Page, int, error) {
	return uc.repo.GetPagesByTag(ctx, tagID, filter)
}
