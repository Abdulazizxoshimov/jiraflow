package page

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/helper"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	pageRepo      repository.PageRepository
	versionRepo   repository.PageVersionRepository
	spaceRepo     repository.SpaceRepository
	issueLinkRepo repository.IssuePageLinkRepository
	issueRepo     repository.IssueRepository
	log           logger.Logger
}

func New(
	pageRepo repository.PageRepository,
	versionRepo repository.PageVersionRepository,
	spaceRepo repository.SpaceRepository,
	issueLinkRepo repository.IssuePageLinkRepository,
	issueRepo repository.IssueRepository,
	log logger.Logger,
) UseCase {
	return &useCase{
		pageRepo:      pageRepo,
		versionRepo:   versionRepo,
		spaceRepo:     spaceRepo,
		issueLinkRepo: issueLinkRepo,
		issueRepo:     issueRepo,
		log:           log,
	}
}

func (uc *useCase) Create(ctx context.Context, spaceID, authorID string, req *entity.CreatePageReq) (*entity.Page, error) {
	isMember, err := uc.spaceRepo.IsMember(ctx, spaceID, authorID)
	if err != nil {
		return nil, fmt.Errorf("page.Create membership check: %w", err)
	}
	if !isMember {
		return nil, apperr.Forbidden("you are not a member of this space")
	}

	pos, err := uc.pageRepo.GetMaxPosition(ctx, spaceID, req.ParentID)
	if err != nil {
		pos = 0
	}

	status := req.Status
	if status == "" {
		status = "draft"
	}

	now := time.Now().UTC()
	p := &entity.Page{
		ID:             uuid.NewString(),
		SpaceID:        spaceID,
		ParentID:       req.ParentID,
		Title:          req.Title,
		Content:        req.Content,
		ContentText:    req.ContentText,
		AuthorID:       authorID,
		LastEditorID:   authorID,
		CurrentVersion: 1,
		Status:         status,
		Position:       pos + 1,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if p.Content == nil {
		p.Content = map[string]any{}
	}

	if err := uc.pageRepo.Create(ctx, p); err != nil {
		uc.log.Error(ctx, "page.Create: db error", logger.String("space_id", spaceID), logger.SafeString("err", err.Error()))
		return nil, err
	}

	v := &entity.PageVersion{
		ID:          uuid.NewString(),
		PageID:      p.ID,
		Version:     1,
		Title:       p.Title,
		Content:     p.Content,
		ContentText: p.ContentText,
		AuthorID:    authorID,
		ChangeNote:  req.ChangeNote,
		CreatedAt:   now,
	}
	_ = uc.versionRepo.Create(ctx, v)

	go uc.processMentions(context.Background(), p.ID, authorID, p.ContentText)

	uc.log.Info(ctx, "page created", logger.String("id", p.ID), logger.String("space_id", spaceID))
	return p, nil
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.Page, error) {
	return uc.pageRepo.GetByID(ctx, id)
}

func (uc *useCase) List(ctx context.Context, filter *entity.PageFilter) ([]*entity.Page, int, error) {
	return uc.pageRepo.List(ctx, filter)
}

func (uc *useCase) Update(ctx context.Context, id, editorID string, req *entity.UpdatePageReq) (*entity.Page, error) {
	p, err := uc.pageRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		p.Title = *req.Title
	}
	if req.Content != nil {
		p.Content = req.Content
	}
	if req.ContentText != nil {
		p.ContentText = *req.ContentText
	}
	if req.Status != nil {
		p.Status = *req.Status
	}
	p.LastEditorID = editorID
	p.CurrentVersion++
	p.UpdatedAt = time.Now().UTC()

	if err := uc.pageRepo.Update(ctx, p); err != nil {
		uc.log.Error(ctx, "page.Update: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return nil, err
	}

	v := &entity.PageVersion{
		ID:          uuid.NewString(),
		PageID:      p.ID,
		Version:     p.CurrentVersion,
		Title:       p.Title,
		Content:     p.Content,
		ContentText: p.ContentText,
		AuthorID:    editorID,
		ChangeNote:  req.ChangeNote,
		CreatedAt:   p.UpdatedAt,
	}
	_ = uc.versionRepo.Create(ctx, v)

	go uc.processMentions(context.Background(), p.ID, editorID, p.ContentText)

	uc.log.Info(ctx, "page updated", logger.String("id", id))
	return p, nil
}

func (uc *useCase) Delete(ctx context.Context, id, actorID string) error {
	p, err := uc.pageRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	allowed := p.AuthorID == actorID
	if !allowed {
		// Space admin ham o'chira oladi
		member, err := uc.spaceRepo.GetMember(ctx, p.SpaceID, actorID)
		if err == nil && member.Role == "admin" {
			allowed = true
		}
	}
	if !allowed {
		uc.log.Warn(ctx, "page.Delete: forbidden", logger.String("id", id), logger.String("actor_id", actorID))
		return apperr.Forbidden("only the author or space admin can delete this page")
	}

	if err := uc.pageRepo.SoftDelete(ctx, id); err != nil {
		uc.log.Error(ctx, "page.Delete: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "page deleted", logger.String("id", id))
	return nil
}

func (uc *useCase) GetTree(ctx context.Context, spaceID string) ([]*entity.PageTree, error) {
	flat, err := uc.pageRepo.GetTree(ctx, spaceID)
	if err != nil {
		return nil, err
	}
	return buildTree(flat), nil
}

func (uc *useCase) Move(ctx context.Context, id string, position int, parentID *string, actorID string) error {
	return uc.pageRepo.UpdatePosition(ctx, id, position, parentID)
}

func (uc *useCase) WatchPage(ctx context.Context, pageID, userID string) error {
	w := &entity.PageWatcher{PageID: pageID, UserID: userID, CreatedAt: time.Now().UTC()}
	return uc.pageRepo.AddWatcher(ctx, w)
}

func (uc *useCase) UnwatchPage(ctx context.Context, pageID, userID string) error {
	return uc.pageRepo.RemoveWatcher(ctx, pageID, userID)
}

func (uc *useCase) ListWatchers(ctx context.Context, pageID string) ([]*entity.PageWatcher, error) {
	return uc.pageRepo.ListWatchers(ctx, pageID)
}

func (uc *useCase) GetWatcherIDs(ctx context.Context, pageID string) ([]string, error) {
	return uc.pageRepo.GetWatcherIDs(ctx, pageID)
}

// processMentions — contentText ichidagi [PROJ-42] mention'larni topib issue-page link yaratadi.
func (uc *useCase) processMentions(ctx context.Context, pageID, editorID, contentText string) {
	if uc.issueLinkRepo == nil || uc.issueRepo == nil {
		return
	}
	keys := helper.ExtractIssueKeys(contentText)
	for _, key := range keys {
		issue, err := uc.issueRepo.GetByKey(ctx, key)
		if err != nil {
			continue
		}
		exists, _ := uc.issueLinkRepo.Exists(ctx, issue.ID, pageID)
		if exists {
			continue
		}
		link := &entity.IssuePageLink{
			ID:        uuid.NewString(),
			IssueID:   issue.ID,
			PageID:    pageID,
			LinkedBy:  editorID,
			CreatedAt: time.Now().UTC(),
		}
		if err := uc.issueLinkRepo.Create(ctx, link); err != nil {
			uc.log.Warn(ctx, "processMentions: link create failed",
				logger.String("key", key), logger.SafeString("err", err.Error()))
		}
	}
}

// ─── Copy ─────────────────────────────────────────────────────────────────────

func (uc *useCase) Copy(ctx context.Context, pageID, actorID string, req *entity.CopyPageReq) (*entity.Page, error) {
	src, err := uc.pageRepo.GetByID(ctx, pageID)
	if err != nil {
		return nil, err
	}

	spaceID := src.SpaceID
	if req.NewSpaceID != nil {
		spaceID = *req.NewSpaceID
	}

	newPage, err := uc.pageRepo.Copy(ctx, pageID, spaceID, req.NewParentID, req.Title, actorID)
	if err != nil {
		return nil, fmt.Errorf("page.Copy: %w", err)
	}

	if req.CopyChildren {
		children, err := uc.pageRepo.GetChildren(ctx, pageID)
		if err != nil {
			return newPage, nil // partial success: return new page, skip children error
		}
		for _, child := range children {
			childReq := &entity.CopyPageReq{
				Title:        child.Title,
				NewSpaceID:   &spaceID,
				NewParentID:  &newPage.ID,
				CopyChildren: true,
			}
			_, _ = uc.Copy(ctx, child.ID, actorID, childReq)
		}
	}

	return newPage, nil
}
