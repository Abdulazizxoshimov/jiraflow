package blueprint

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, req *entity.CreateBlueprintReq) (*entity.Blueprint, error)
	GetByID(ctx context.Context, id string) (*entity.Blueprint, error)
	List(ctx context.Context) ([]*entity.Blueprint, error)
	Delete(ctx context.Context, id string) error
	CreatePage(ctx context.Context, blueprintID, actorID string, req *entity.CreatePageFromBlueprintReq) (*entity.Page, error)
}
