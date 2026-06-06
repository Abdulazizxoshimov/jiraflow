package label

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, projectID string, req *entity.CreateLabelReq) (*entity.Label, error)
	GetByID(ctx context.Context, id string) (*entity.Label, error)
	ListByProject(ctx context.Context, projectID string) ([]*entity.Label, error)
	Update(ctx context.Context, id string, req *entity.UpdateLabelReq) (*entity.Label, error)
	Delete(ctx context.Context, id string) error
}
