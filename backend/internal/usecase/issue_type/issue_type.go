package issue_type

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
)

type useCase struct {
	repo repository.IssueTypeRepository
}

func New(repo repository.IssueTypeRepository) UseCase {
	return &useCase{repo: repo}
}

func (uc *useCase) CreateType(ctx context.Context, req *entity.CreateIssueTypeReq) (*entity.IssueType, error) {
	t := &entity.IssueType{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Description: req.Description,
		IconURL:     req.IconURL,
		Color:       req.Color,
		IsSubtask:   req.IsSubtask,
		IsSystem:    false,
	}
	if err := uc.repo.CreateType(ctx, t); err != nil {
		return nil, fmt.Errorf("issue_type.CreateType: %w", err)
	}
	return t, nil
}

func (uc *useCase) ListTypes(ctx context.Context) ([]*entity.IssueType, error) {
	return uc.repo.ListTypes(ctx)
}

func (uc *useCase) GetTypeByID(ctx context.Context, id string) (*entity.IssueType, error) {
	return uc.repo.GetTypeByID(ctx, id)
}

func (uc *useCase) DeleteType(ctx context.Context, id string) error {
	return uc.repo.DeleteType(ctx, id)
}

func (uc *useCase) CreateScheme(ctx context.Context, req *entity.CreateIssueTypeSchemeReq) (*entity.IssueTypeScheme, error) {
	s := &entity.IssueTypeScheme{
		ID:        uuid.NewString(),
		Name:      req.Name,
		ProjectID: req.ProjectID,
	}
	if err := uc.repo.CreateScheme(ctx, s, req.IssueTypeIDs); err != nil {
		return nil, fmt.Errorf("issue_type.CreateScheme: %w", err)
	}
	return uc.repo.GetSchemeByID(ctx, s.ID)
}

func (uc *useCase) GetSchemeByID(ctx context.Context, id string) (*entity.IssueTypeScheme, error) {
	return uc.repo.GetSchemeByID(ctx, id)
}

func (uc *useCase) GetSchemeByProject(ctx context.Context, projectID string) (*entity.IssueTypeScheme, error) {
	return uc.repo.GetSchemeByProject(ctx, projectID)
}

func (uc *useCase) ListSchemes(ctx context.Context) ([]*entity.IssueTypeScheme, error) {
	return uc.repo.ListSchemes(ctx)
}

func (uc *useCase) DeleteScheme(ctx context.Context, id string) error {
	return uc.repo.DeleteScheme(ctx, id)
}
