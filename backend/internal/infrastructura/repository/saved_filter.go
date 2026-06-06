package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type SavedFilterRepository interface {
	Create(ctx context.Context, userID string, req *entity.CreateSavedFilterReq) (*entity.SavedFilter, error)
	GetByID(ctx context.Context, id string) (*entity.SavedFilter, error)
	List(ctx context.Context, userID, filterType string) ([]*entity.SavedFilter, error)
	Update(ctx context.Context, id string, req *entity.UpdateSavedFilterReq) (*entity.SavedFilter, error)
	Delete(ctx context.Context, id, userID string) error
}
