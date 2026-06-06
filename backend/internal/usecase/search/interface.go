package search

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Search(ctx context.Context, filter *entity.SearchFilter) ([]*entity.SearchResult, int, error)
	Suggest(ctx context.Context, query string, limit int) ([]*entity.SearchSuggestion, error)
}
