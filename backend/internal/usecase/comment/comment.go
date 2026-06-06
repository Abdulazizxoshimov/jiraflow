package comment

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/tiptap"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/notification"
)

type useCase struct {
	repo       repository.CommentRepository
	issueRepo  repository.IssueRepository
	pageRepo   repository.PageRepository
	dispatcher notification.Dispatcher
	log        logger.Logger
}

func New(
	repo repository.CommentRepository,
	issueRepo repository.IssueRepository,
	pageRepo repository.PageRepository,
	dispatcher notification.Dispatcher,
	log logger.Logger,
) UseCase {
	return &useCase{repo: repo, issueRepo: issueRepo, pageRepo: pageRepo, dispatcher: dispatcher, log: log}
}

func (uc *useCase) Create(ctx context.Context, parentType, parentID, authorID string, req *entity.CreateCommentReq) (*entity.Comment, error) {
	now := time.Now().UTC()
	c := &entity.Comment{
		ID:          uuid.NewString(),
		ParentType:  parentType,
		ParentID:    parentID,
		AuthorID:    authorID,
		Content:     req.Content,
		ContentText: req.ContentText,
		ReplyToID:   req.ReplyToID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := uc.repo.Create(ctx, c); err != nil {
		uc.log.Error(ctx, "comment.Create: db error", logger.String("parent_id", parentID), logger.SafeString("err", err.Error()))
		return nil, err
	}

	go func() {
		switch parentType {
		case "issues":
			watchers, _ := uc.issueRepo.ListWatchers(ctx, parentID)
			ids := make([]string, 0, len(watchers))
			for _, w := range watchers {
				ids = append(ids, w.UserID)
			}
			uc.dispatcher.IssueCommented(ctx, parentID, ids, authorID)
		case "pages":
			ids, _ := uc.pageRepo.GetWatcherIDs(ctx, parentID)
			uc.dispatcher.PageCommented(ctx, parentID, ids, authorID)
			if mentionIDs := tiptap.ExtractMentionIDs(c.Content); len(mentionIDs) > 0 {
				uc.dispatcher.PageMentioned(ctx, parentID, mentionIDs, authorID)
			}
		}
	}()

	uc.log.Info(ctx, "comment created", logger.String("id", c.ID), logger.String("parent_id", parentID))
	return c, nil
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.Comment, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *useCase) ListByParent(ctx context.Context, parentType, parentID string, filter *entity.Filter) ([]*entity.Comment, int, error) {
	return uc.repo.ListByParent(ctx, parentType, parentID, filter)
}

func (uc *useCase) Update(ctx context.Context, id, actorID string, req *entity.UpdateCommentReq) (*entity.Comment, error) {
	c, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c.AuthorID != actorID {
		uc.log.Warn(ctx, "comment.Update: forbidden", logger.String("id", id), logger.String("actor_id", actorID))
		return nil, apperr.Forbidden("only the author can edit this comment")
	}
	c.Content = req.Content
	c.ContentText = req.ContentText
	if err := uc.repo.Update(ctx, c); err != nil {
		uc.log.Error(ctx, "comment.Update: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return nil, err
	}
	_ = uc.repo.DeleteMentions(ctx, id)
	uc.log.Info(ctx, "comment updated", logger.String("id", id))
	return c, nil
}

func (uc *useCase) Delete(ctx context.Context, id, actorID string) error {
	c, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if c.AuthorID != actorID {
		uc.log.Warn(ctx, "comment.Delete: forbidden", logger.String("id", id), logger.String("actor_id", actorID))
		return apperr.Forbidden("only the author can delete this comment")
	}
	if err := uc.repo.SoftDelete(ctx, id); err != nil {
		uc.log.Error(ctx, "comment.Delete: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "comment deleted", logger.String("id", id))
	return nil
}

func (uc *useCase) ToggleReaction(ctx context.Context, commentID, userID, emoji string) error {
	if _, err := uc.repo.GetByID(ctx, commentID); err != nil {
		return err
	}
	return uc.repo.ToggleReaction(ctx, commentID, userID, emoji)
}

func (uc *useCase) ListReactions(ctx context.Context, commentID, viewerID string) ([]entity.CommentReactionSummary, error) {
	return uc.repo.ListReactions(ctx, commentID, viewerID)
}
