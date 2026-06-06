package page_template

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, spaceID *string, createdBy string, req *entity.CreatePageTemplateReq) (*entity.PageTemplate, error)
	GetByID(ctx context.Context, id string) (*entity.PageTemplate, error)
	List(ctx context.Context, filter *entity.PageTemplateFilter) ([]*entity.PageTemplate, int, error)
	Update(ctx context.Context, id, actorID string, req *entity.UpdatePageTemplateReq) (*entity.PageTemplate, error)
	Delete(ctx context.Context, id, actorID string) error
}
