package page_macro

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Upsert(ctx context.Context, pageID string, req *entity.UpsertPageMacroReq) (*entity.PageMacro, error)
	ListByPage(ctx context.Context, pageID string) ([]*entity.PageMacro, error)
	GetByID(ctx context.Context, id string) (*entity.PageMacro, error)
	Delete(ctx context.Context, id string) error
}
