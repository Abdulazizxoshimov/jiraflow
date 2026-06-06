package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type SpaceExportRepository interface {
	Create(ctx context.Context, e *entity.SpaceExport) error
	GetByID(ctx context.Context, id string) (*entity.SpaceExport, error)
	UpdateStatus(ctx context.Context, id, status string, fileURL, errorMsg *string) error
	ListBySpace(ctx context.Context, spaceID string) ([]*entity.SpaceExport, error)
}
