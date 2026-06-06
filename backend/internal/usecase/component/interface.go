package component

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, projectID string, req *entity.CreateComponentReq) (*entity.Component, error)
	GetByID(ctx context.Context, id string) (*entity.Component, error)
	List(ctx context.Context, projectID string) ([]*entity.Component, error)
	Update(ctx context.Context, id string, req *entity.UpdateComponentReq) (*entity.Component, error)
	Delete(ctx context.Context, id string) error
}
