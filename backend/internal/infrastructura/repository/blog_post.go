package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type BlogPostRepository interface {
	Create(ctx context.Context, spaceID, authorID string, req *entity.CreateBlogPostReq) (*entity.BlogPost, error)
	GetByID(ctx context.Context, id string) (*entity.BlogPost, error)
	List(ctx context.Context, filter entity.ListBlogPostsFilter) ([]*entity.BlogPost, int, error)
	Update(ctx context.Context, id string, req *entity.UpdateBlogPostReq) (*entity.BlogPost, error)
	Delete(ctx context.Context, id string) error
	Publish(ctx context.Context, id string) error
	Unpublish(ctx context.Context, id string) error
}
