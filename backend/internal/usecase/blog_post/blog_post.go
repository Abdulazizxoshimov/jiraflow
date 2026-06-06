package blog_post

import (
	"context"
	"fmt"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
)

type useCase struct {
	repo repository.BlogPostRepository
}

func New(repo repository.BlogPostRepository) UseCase {
	return &useCase{repo: repo}
}

func (uc *useCase) Create(ctx context.Context, spaceID, actorID string, req *entity.CreateBlogPostReq) (*entity.BlogPost, error) {
	return uc.repo.Create(ctx, spaceID, actorID, req)
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.BlogPost, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *useCase) List(ctx context.Context, filter entity.ListBlogPostsFilter) ([]*entity.BlogPost, int, error) {
	return uc.repo.List(ctx, filter)
}

func (uc *useCase) Update(ctx context.Context, id, actorID string, req *entity.UpdateBlogPostReq) (*entity.BlogPost, error) {
	bp, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if bp.AuthorID != actorID {
		return nil, apperr.Forbidden("only the author can edit this blog post")
	}
	return uc.repo.Update(ctx, id, req)
}

func (uc *useCase) Delete(ctx context.Context, id, actorID string) error {
	bp, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if bp.AuthorID != actorID {
		return apperr.Forbidden("only the author can delete this blog post")
	}
	return uc.repo.Delete(ctx, id)
}

func (uc *useCase) Publish(ctx context.Context, id, actorID string) error {
	bp, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if bp.AuthorID != actorID {
		return fmt.Errorf("blog_post.Publish: %w", apperr.Forbidden("only the author can publish"))
	}
	return uc.repo.Publish(ctx, id)
}

func (uc *useCase) Unpublish(ctx context.Context, id, actorID string) error {
	bp, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if bp.AuthorID != actorID {
		return fmt.Errorf("blog_post.Unpublish: %w", apperr.Forbidden("only the author can unpublish"))
	}
	return uc.repo.Unpublish(ctx, id)
}
