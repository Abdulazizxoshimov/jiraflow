package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type PageRepository interface {
	Create(ctx context.Context, p *entity.Page) error
	GetByID(ctx context.Context, id string) (*entity.Page, error)
	List(ctx context.Context, filter *entity.PageFilter) ([]*entity.Page, int, error)
	Update(ctx context.Context, p *entity.Page) error
	SoftDelete(ctx context.Context, id string) error
	GetTree(ctx context.Context, spaceID string) ([]*entity.PageTree, error)
	UpdatePosition(ctx context.Context, id string, position int, parentID *string) error
	GetMaxPosition(ctx context.Context, spaceID string, parentID *string) (int, error)

	AddWatcher(ctx context.Context, w *entity.PageWatcher) error
	RemoveWatcher(ctx context.Context, pageID, userID string) error
	ListWatchers(ctx context.Context, pageID string) ([]*entity.PageWatcher, error)
	IsWatcher(ctx context.Context, pageID, userID string) (bool, error)
	GetWatcherIDs(ctx context.Context, pageID string) ([]string, error)

	Copy(ctx context.Context, srcID, newSpaceID string, newParentID *string, newTitle string, authorID string) (*entity.Page, error)
	GetChildren(ctx context.Context, parentID string) ([]*entity.Page, error)
}
