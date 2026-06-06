package inline_comment

import (
	"context"
	"fmt"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	repo repository.InlineCommentRepository
	log  logger.Logger
}

func New(repo repository.InlineCommentRepository, log logger.Logger) UseCase {
	return &useCase{repo: repo, log: log}
}

func (uc *useCase) Create(ctx context.Context, pageID, authorID string, req *entity.CreateInlineCommentReq) (*entity.InlineComment, error) {
	c := &entity.InlineComment{
		PageID:    pageID,
		AuthorID:  authorID,
		AnchorID:  req.AnchorID,
		QuoteText: req.QuoteText,
		Body:      req.Body,
	}
	if err := uc.repo.Create(ctx, c); err != nil {
		uc.log.Error(ctx, "inlineComment.Create: db error", logger.String("page_id", pageID), logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("inlineComment.Create: %w", err)
	}
	uc.log.Info(ctx, "inline comment created", logger.String("id", c.ID))
	return c, nil
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.InlineComment, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *useCase) ListByPage(ctx context.Context, pageID string) ([]*entity.InlineComment, error) {
	return uc.repo.ListByPage(ctx, pageID)
}

func (uc *useCase) ListByAnchor(ctx context.Context, pageID, anchorID string) ([]*entity.InlineComment, error) {
	return uc.repo.ListByAnchor(ctx, pageID, anchorID)
}

func (uc *useCase) Update(ctx context.Context, id, actorID string, req *entity.UpdateInlineCommentReq) (*entity.InlineComment, error) {
	c, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c.AuthorID != actorID {
		return nil, apperr.Forbidden("only the author can edit this comment")
	}
	c.Body = req.Body
	if err := uc.repo.Update(ctx, c); err != nil {
		uc.log.Error(ctx, "inlineComment.Update: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("inlineComment.Update: %w", err)
	}
	return c, nil
}

func (uc *useCase) Resolve(ctx context.Context, id, resolverID string) error {
	if _, err := uc.repo.GetByID(ctx, id); err != nil {
		return err
	}
	return uc.repo.Resolve(ctx, id, resolverID)
}

func (uc *useCase) Unresolve(ctx context.Context, id, actorID string) error {
	if _, err := uc.repo.GetByID(ctx, id); err != nil {
		return err
	}
	return uc.repo.Unresolve(ctx, id)
}

func (uc *useCase) Delete(ctx context.Context, id, actorID string) error {
	c, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if c.AuthorID != actorID {
		return apperr.Forbidden("only the author can delete this comment")
	}
	if err := uc.repo.SoftDelete(ctx, id); err != nil {
		uc.log.Error(ctx, "inlineComment.Delete: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return fmt.Errorf("inlineComment.Delete: %w", err)
	}
	return nil
}
