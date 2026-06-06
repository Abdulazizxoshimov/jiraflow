package space_export

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	RequestExport(ctx context.Context, spaceID, requestedBy string) (*entity.SpaceExport, error)
	GetExport(ctx context.Context, id string) (*entity.SpaceExport, error)
	ListExports(ctx context.Context, spaceID string) ([]*entity.SpaceExport, error)
}
