package space_category

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, req *entity.CreateSpaceCategoryReq) (*entity.SpaceCategory, error)
	GetByID(ctx context.Context, id string) (*entity.SpaceCategory, error)
	List(ctx context.Context) ([]*entity.SpaceCategory, error)
	Update(ctx context.Context, id string, req *entity.UpdateSpaceCategoryReq) (*entity.SpaceCategory, error)
	Delete(ctx context.Context, id string) error
}
