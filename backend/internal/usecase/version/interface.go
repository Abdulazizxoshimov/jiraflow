package version

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, projectID string, req *entity.CreateVersionReq) (*entity.Version, error)
	GetByID(ctx context.Context, id string) (*entity.Version, error)
	List(ctx context.Context, projectID string) ([]*entity.Version, error)
	Update(ctx context.Context, id string, req *entity.UpdateVersionReq) (*entity.Version, error)
	Release(ctx context.Context, id string, req *entity.ReleaseVersionReq) (*entity.Version, error)
	Archive(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}
