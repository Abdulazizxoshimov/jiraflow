package favorite

import (
	"context"
	"fmt"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	repo      repository.FavoriteRepository
	pageRepo  repository.PageRepository
	spaceRepo repository.SpaceRepository
	log       logger.Logger
}

func New(repo repository.FavoriteRepository, pageRepo repository.PageRepository, spaceRepo repository.SpaceRepository, log logger.Logger) UseCase {
	return &useCase{repo: repo, pageRepo: pageRepo, spaceRepo: spaceRepo, log: log}
}

func (uc *useCase) Add(ctx context.Context, userID string, req *entity.AddFavoriteReq) (*entity.Favorite, error) {
	// Validate that entity exists
	switch req.EntityType {
	case "page":
		if _, err := uc.pageRepo.GetByID(ctx, req.EntityID); err != nil {
			return nil, err
		}
	case "space":
		if _, err := uc.spaceRepo.GetByID(ctx, req.EntityID); err != nil {
			return nil, err
		}
	default:
		return nil, apperr.BadRequest("invalid entity_type")
	}

	fav := &entity.Favorite{
		UserID:     userID,
		EntityType: req.EntityType,
		EntityID:   req.EntityID,
	}
	if err := uc.repo.Add(ctx, fav); err != nil {
		uc.log.Error(ctx, "favorite.Add: db error", logger.String("user_id", userID), logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("favorite.Add: %w", err)
	}
	uc.log.Info(ctx, "favorite added", logger.String("user_id", userID), logger.String("entity_id", req.EntityID))
	return fav, nil
}

func (uc *useCase) Remove(ctx context.Context, userID, entityType, entityID string) error {
	if err := uc.repo.Remove(ctx, userID, entityType, entityID); err != nil {
		uc.log.Error(ctx, "favorite.Remove: db error", logger.String("user_id", userID), logger.SafeString("err", err.Error()))
		return fmt.Errorf("favorite.Remove: %w", err)
	}
	return nil
}

func (uc *useCase) List(ctx context.Context, userID string, filter *entity.FavoriteFilter) ([]*entity.Favorite, int, error) {
	favs, total, err := uc.repo.List(ctx, userID, filter)
	if err != nil {
		return nil, 0, err
	}

	for _, fav := range favs {
		switch fav.EntityType {
		case "page":
			page, err := uc.pageRepo.GetByID(ctx, fav.EntityID)
			if err == nil {
				fav.Page = page
			}
		case "space":
			space, err := uc.spaceRepo.GetByID(ctx, fav.EntityID)
			if err == nil {
				fav.Space = space
			}
		}
	}
	return favs, total, nil
}

func (uc *useCase) IsFavorite(ctx context.Context, userID, entityType, entityID string) (bool, error) {
	return uc.repo.Exists(ctx, userID, entityType, entityID)
}
