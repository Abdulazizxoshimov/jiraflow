package permission_scheme

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, req *entity.CreatePermissionSchemeReq, createdBy string) (*entity.PermissionScheme, error)
	Get(ctx context.Context, id string) (*entity.PermissionScheme, error)
	List(ctx context.Context) ([]*entity.PermissionScheme, error)
	Update(ctx context.Context, id string, req *entity.UpdatePermissionSchemeReq) (*entity.PermissionScheme, error)
	Delete(ctx context.Context, id string) error

	AddGrant(ctx context.Context, schemeID string, req *entity.AddGrantReq) (*entity.PermissionSchemeGrant, error)
	RemoveGrant(ctx context.Context, schemeID, grantID string) error

	AssignToProject(ctx context.Context, projectID, schemeID string) error
	GetByProject(ctx context.Context, projectID string) (*entity.PermissionScheme, error)
}
