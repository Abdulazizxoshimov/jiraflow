package security_scheme

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
)

type UseCase interface {
	Create(ctx context.Context, req *entity.CreateSecuritySchemeReq) (*entity.SecurityScheme, error)
	GetByID(ctx context.Context, id string) (*entity.SecurityScheme, error)
	List(ctx context.Context, projectID string) ([]*entity.SecurityScheme, error)
	Delete(ctx context.Context, id string) error
	AddLevel(ctx context.Context, schemeID string, req *entity.CreateSecurityLevelReq) (*entity.SecurityLevel, error)
	GetLevel(ctx context.Context, levelID string) (*entity.SecurityLevel, error)
	DeleteLevel(ctx context.Context, levelID string) error
	AddMember(ctx context.Context, levelID string, req *entity.CreateSecurityLevelMemberReq) (*entity.SecurityLevelMember, error)
	DeleteMember(ctx context.Context, memberID string) error
}

type useCase struct {
	repo repository.SecuritySchemeRepository
}

func New(repo repository.SecuritySchemeRepository) UseCase {
	return &useCase{repo: repo}
}

func (uc *useCase) Create(ctx context.Context, req *entity.CreateSecuritySchemeReq) (*entity.SecurityScheme, error) {
	return uc.repo.Create(ctx, req)
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.SecurityScheme, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *useCase) List(ctx context.Context, projectID string) ([]*entity.SecurityScheme, error) {
	return uc.repo.List(ctx, projectID)
}

func (uc *useCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *useCase) AddLevel(ctx context.Context, schemeID string, req *entity.CreateSecurityLevelReq) (*entity.SecurityLevel, error) {
	return uc.repo.AddLevel(ctx, schemeID, req)
}

func (uc *useCase) GetLevel(ctx context.Context, levelID string) (*entity.SecurityLevel, error) {
	return uc.repo.GetLevel(ctx, levelID)
}

func (uc *useCase) DeleteLevel(ctx context.Context, levelID string) error {
	return uc.repo.DeleteLevel(ctx, levelID)
}

func (uc *useCase) AddMember(ctx context.Context, levelID string, req *entity.CreateSecurityLevelMemberReq) (*entity.SecurityLevelMember, error) {
	return uc.repo.AddMember(ctx, levelID, req)
}

func (uc *useCase) DeleteMember(ctx context.Context, memberID string) error {
	return uc.repo.DeleteMember(ctx, memberID)
}
