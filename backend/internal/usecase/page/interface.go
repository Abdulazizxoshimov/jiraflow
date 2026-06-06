package page

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, spaceID, authorID string, req *entity.CreatePageReq) (*entity.Page, error)
	GetByID(ctx context.Context, id string) (*entity.Page, error)
	List(ctx context.Context, filter *entity.PageFilter) ([]*entity.Page, int, error)
	Update(ctx context.Context, id, editorID string, req *entity.UpdatePageReq) (*entity.Page, error)
	Delete(ctx context.Context, id, actorID string) error
	GetTree(ctx context.Context, spaceID string) ([]*entity.PageTree, error)
	Move(ctx context.Context, id string, position int, parentID *string, actorID string) error

	WatchPage(ctx context.Context, pageID, userID string) error
	UnwatchPage(ctx context.Context, pageID, userID string) error
	ListWatchers(ctx context.Context, pageID string) ([]*entity.PageWatcher, error)
	GetWatcherIDs(ctx context.Context, pageID string) ([]string, error)

	Copy(ctx context.Context, pageID, actorID string, req *entity.CopyPageReq) (*entity.Page, error)
}
