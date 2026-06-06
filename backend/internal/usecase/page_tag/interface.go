package page_tag

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, spaceID string, req *entity.CreatePageTagReq) (*entity.PageTag, error)
	GetByID(ctx context.Context, id string) (*entity.PageTag, error)
	List(ctx context.Context, spaceID string) ([]*entity.PageTag, error)
	Update(ctx context.Context, id string, req *entity.UpdatePageTagReq) (*entity.PageTag, error)
	Delete(ctx context.Context, id string) error

	SetPageTags(ctx context.Context, pageID string, req *entity.SetPageTagsReq) error
	GetPageTags(ctx context.Context, pageID string) ([]*entity.PageTag, error)
	GetPagesByTag(ctx context.Context, tagID string, filter *entity.Filter) ([]*entity.Page, int, error)
}
