package data_import

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	ImportJira(ctx context.Context, userID string, data []byte) (*entity.DataImport, error)
	ImportTrello(ctx context.Context, userID string, data []byte) (*entity.DataImport, error)
	ImportLinear(ctx context.Context, userID string, data []byte) (*entity.DataImport, error)
	GetStatus(ctx context.Context, id string) (*entity.DataImport, error)
}
