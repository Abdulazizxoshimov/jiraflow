package blog_post

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, spaceID, actorID string, req *entity.CreateBlogPostReq) (*entity.BlogPost, error)
	GetByID(ctx context.Context, id string) (*entity.BlogPost, error)
	List(ctx context.Context, filter entity.ListBlogPostsFilter) ([]*entity.BlogPost, int, error)
	Update(ctx context.Context, id, actorID string, req *entity.UpdateBlogPostReq) (*entity.BlogPost, error)
	Delete(ctx context.Context, id, actorID string) error
	Publish(ctx context.Context, id, actorID string) error
	Unpublish(ctx context.Context, id, actorID string) error
}
