package saved_filter

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
)

type useCase struct {
	repo repository.SavedFilterRepository
}

func New(repo repository.SavedFilterRepository) UseCase {
	return &useCase{repo: repo}
}

func (uc *useCase) Create(ctx context.Context, actorID string, req *entity.CreateSavedFilterReq) (*entity.SavedFilter, error) {
	return uc.repo.Create(ctx, actorID, req)
}

func (uc *useCase) GetByID(ctx context.Context, id, actorID string) (*entity.SavedFilter, error) {
	sf, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if sf.UserID != actorID && !sf.IsShared {
		return nil, apperr.Forbidden("access denied")
	}
	return sf, nil
}

func (uc *useCase) List(ctx context.Context, actorID, filterType string) ([]*entity.SavedFilter, error) {
	return uc.repo.List(ctx, actorID, filterType)
}

func (uc *useCase) Update(ctx context.Context, id, actorID string, req *entity.UpdateSavedFilterReq) (*entity.SavedFilter, error) {
	sf, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if sf.UserID != actorID {
		return nil, apperr.Forbidden("only the owner can edit this filter")
	}
	return uc.repo.Update(ctx, id, req)
}

func (uc *useCase) Delete(ctx context.Context, id, actorID string) error {
	sf, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if sf.UserID != actorID {
		return apperr.Forbidden("only the owner can delete this filter")
	}
	return uc.repo.Delete(ctx, id, actorID)
}
