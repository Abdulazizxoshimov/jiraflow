package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type PageRestrictionRepository interface {
	Set(ctx context.Context, pageID string, items []entity.PageRestrictionItem) error
	List(ctx context.Context, pageID string) ([]*entity.PageRestriction, error)
	Clear(ctx context.Context, pageID string) error
	CanAccess(ctx context.Context, pageID, userID, accessType string) (bool, error)
}
