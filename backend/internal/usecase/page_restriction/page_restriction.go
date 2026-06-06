package page_restriction

import (
	"context"
	"fmt"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	repo      repository.PageRestrictionRepository
	pageRepo  repository.PageRepository
	spaceRepo repository.SpaceRepository
	log       logger.Logger
}

func New(repo repository.PageRestrictionRepository, pageRepo repository.PageRepository, spaceRepo repository.SpaceRepository, log logger.Logger) UseCase {
	return &useCase{repo: repo, pageRepo: pageRepo, spaceRepo: spaceRepo, log: log}
}

func (uc *useCase) Set(ctx context.Context, pageID, actorID string, req *entity.SetPageRestrictionsReq) error {
	page, err := uc.pageRepo.GetByID(ctx, pageID)
	if err != nil {
		return err
	}
	if page.AuthorID != actorID {
		member, err := uc.spaceRepo.GetMember(ctx, page.SpaceID, actorID)
		if err != nil || member.Role != "admin" {
			return fmt.Errorf("page_restriction.Set: only author or space admin can set restrictions")
		}
	}
	if err := uc.repo.Set(ctx, pageID, req.Restrictions); err != nil {
		uc.log.Error(ctx, "pageRestriction.Set: db error", logger.String("page_id", pageID), logger.SafeString("err", err.Error()))
		return fmt.Errorf("pageRestriction.Set: %w", err)
	}
	return nil
}

func (uc *useCase) List(ctx context.Context, pageID string) ([]*entity.PageRestriction, error) {
	return uc.repo.List(ctx, pageID)
}

func (uc *useCase) Clear(ctx context.Context, pageID, actorID string) error {
	page, err := uc.pageRepo.GetByID(ctx, pageID)
	if err != nil {
		return err
	}
	if page.AuthorID != actorID {
		member, err := uc.spaceRepo.GetMember(ctx, page.SpaceID, actorID)
		if err != nil || member.Role != "admin" {
			return fmt.Errorf("page_restriction.Clear: only author or space admin can clear restrictions")
		}
	}
	return uc.repo.Clear(ctx, pageID)
}

func (uc *useCase) CheckAccess(ctx context.Context, pageID, userID, accessType string) (*entity.PageAccessInfo, error) {
	canView, err := uc.repo.CanAccess(ctx, pageID, userID, "view")
	if err != nil {
		return nil, fmt.Errorf("pageRestriction.CheckAccess view: %w", err)
	}
	canEdit, err := uc.repo.CanAccess(ctx, pageID, userID, "edit")
	if err != nil {
		return nil, fmt.Errorf("pageRestriction.CheckAccess edit: %w", err)
	}
	return &entity.PageAccessInfo{CanView: canView, CanEdit: canEdit}, nil
}
