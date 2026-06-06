package search

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	searchRepo repository.SearchRepository
	log        logger.Logger
}

func New(searchRepo repository.SearchRepository, log logger.Logger) UseCase {
	return &useCase{searchRepo: searchRepo, log: log}
}

func (uc *useCase) Search(ctx context.Context, filter *entity.SearchFilter) ([]*entity.SearchResult, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	results, total, err := uc.searchRepo.Search(ctx, filter)
	if err != nil {
		uc.log.Error(ctx, "search.Search: db error", logger.SafeString("err", err.Error()))
		return nil, 0, err
	}
	uc.log.Debug(ctx, "search executed", logger.Int("total", total))
	return results, total, nil
}

func (uc *useCase) Suggest(ctx context.Context, query string, limit int) ([]*entity.SearchSuggestion, error) {
	if len(query) < 2 {
		return nil, nil
	}
	suggestions, err := uc.searchRepo.Suggest(ctx, query, limit)
	if err != nil {
		uc.log.Error(ctx, "search.Suggest: db error", logger.SafeString("err", err.Error()))
		return nil, err
	}
	return suggestions, nil
}
