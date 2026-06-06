package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type PermissionSchemeRepository interface {
	Create(ctx context.Context, scheme *entity.PermissionScheme) error
	GetByID(ctx context.Context, id string) (*entity.PermissionScheme, error)
	List(ctx context.Context) ([]*entity.PermissionScheme, error)
	Update(ctx context.Context, scheme *entity.PermissionScheme) error
	Delete(ctx context.Context, id string) error

	AddGrant(ctx context.Context, grant *entity.PermissionSchemeGrant) error
	RemoveGrant(ctx context.Context, grantID string) error
	ListGrants(ctx context.Context, schemeID string) ([]*entity.PermissionSchemeGrant, error)

	AssignToProject(ctx context.Context, projectID, schemeID string) error
	GetByProject(ctx context.Context, projectID string) (*entity.PermissionScheme, error)
}
