package permission_scheme

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
)

type useCase struct {
	repo repository.PermissionSchemeRepository
}

func New(repo repository.PermissionSchemeRepository) UseCase {
	return &useCase{repo: repo}
}

func (uc *useCase) Create(ctx context.Context, req *entity.CreatePermissionSchemeReq, createdBy string) (*entity.PermissionScheme, error) {
	now := time.Now().UTC()
	s := &entity.PermissionScheme{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := uc.repo.Create(ctx, s); err != nil {
		return nil, fmt.Errorf("permission_scheme.Create: %w", err)
	}
	return s, nil
}

func (uc *useCase) Get(ctx context.Context, id string) (*entity.PermissionScheme, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *useCase) List(ctx context.Context) ([]*entity.PermissionScheme, error) {
	return uc.repo.List(ctx)
}

func (uc *useCase) Update(ctx context.Context, id string, req *entity.UpdatePermissionSchemeReq) (*entity.PermissionScheme, error) {
	s, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		s.Name = *req.Name
	}
	if req.Description != nil {
		s.Description = *req.Description
	}
	s.UpdatedAt = time.Now().UTC()
	if err := uc.repo.Update(ctx, s); err != nil {
		return nil, fmt.Errorf("permission_scheme.Update: %w", err)
	}
	return s, nil
}

func (uc *useCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *useCase) AddGrant(ctx context.Context, schemeID string, req *entity.AddGrantReq) (*entity.PermissionSchemeGrant, error) {
	g := &entity.PermissionSchemeGrant{
		ID:         uuid.NewString(),
		SchemeID:   schemeID,
		Permission: req.Permission,
		HolderType: req.HolderType,
		HolderID:   req.HolderID,
		CreatedAt:  time.Now().UTC(),
	}
	if err := uc.repo.AddGrant(ctx, g); err != nil {
		return nil, fmt.Errorf("permission_scheme.AddGrant: %w", err)
	}
	return g, nil
}

func (uc *useCase) RemoveGrant(ctx context.Context, schemeID, grantID string) error {
	return uc.repo.RemoveGrant(ctx, grantID)
}

func (uc *useCase) AssignToProject(ctx context.Context, projectID, schemeID string) error {
	return uc.repo.AssignToProject(ctx, projectID, schemeID)
}

func (uc *useCase) GetByProject(ctx context.Context, projectID string) (*entity.PermissionScheme, error) {
	return uc.repo.GetByProject(ctx, projectID)
}
