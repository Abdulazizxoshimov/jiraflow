package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type DataImportRepository interface {
	Create(ctx context.Context, imp *entity.DataImport) error
	GetByID(ctx context.Context, id string) (*entity.DataImport, error)
	UpdateStatus(ctx context.Context, id, status string, total, processed int, errMsg string) error
	MarkCompleted(ctx context.Context, id string) error
}
