package page_restriction

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Set(ctx context.Context, pageID, actorID string, req *entity.SetPageRestrictionsReq) error
	List(ctx context.Context, pageID string) ([]*entity.PageRestriction, error)
	Clear(ctx context.Context, pageID, actorID string) error
	CheckAccess(ctx context.Context, pageID, userID, accessType string) (*entity.PageAccessInfo, error)
}
