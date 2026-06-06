package saved_filter

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, actorID string, req *entity.CreateSavedFilterReq) (*entity.SavedFilter, error)
	GetByID(ctx context.Context, id, actorID string) (*entity.SavedFilter, error)
	List(ctx context.Context, actorID, filterType string) ([]*entity.SavedFilter, error)
	Update(ctx context.Context, id, actorID string, req *entity.UpdateSavedFilterReq) (*entity.SavedFilter, error)
	Delete(ctx context.Context, id, actorID string) error
}
