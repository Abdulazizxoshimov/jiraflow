package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type searchRepo struct {
	db *pgxpool.Pool
}

func NewSearchRepo(p *pg.Postgres) repository.SearchRepository {
	return &searchRepo{db: p.DB}
}

func (r *searchRepo) Search(ctx context.Context, filter *entity.SearchFilter) ([]*entity.SearchResult, int, error) {
	types := filter.Types
	if len(types) == 0 {
		types = []entity.SearchResultType{
			entity.SearchResultIssue,
			entity.SearchResultPage,
			entity.SearchResultProject,
			entity.SearchResultSpace,
		}
	}

	var parts []string
	var args []any
	argIdx := 1

	tsquery := strings.TrimSpace(filter.Query)

	for _, t := range types {
		switch t {
		case entity.SearchResultIssue:
			cond := ""
			if filter.ProjectID != "" {
				cond = fmt.Sprintf(" AND i.project_id=$%d", argIdx+1)
				args = append(args, tsquery, filter.ProjectID)
				argIdx += 2
			} else {
				args = append(args, tsquery)
				argIdx++
			}
			parts = append(parts, fmt.Sprintf(`
				SELECT 'issue' AS type, i.id, i.title,
				       LEFT(COALESCE(i.description,''), 200) AS excerpt,
				       i.updated_at,
				       COALESCE(ts_rank(i.search_vector, plainto_tsquery('english', $%d)), 0) AS score
				FROM issues i
				WHERE i.deleted_at IS NULL
				  AND (i.search_vector @@ plainto_tsquery('english', $%d)
				       OR i.title ILIKE '%%' || $%d || '%%')
				  %s
			`, argIdx-1, argIdx-1, argIdx-1, cond))

		case entity.SearchResultPage:
			cond := ""
			if filter.SpaceID != "" {
				cond = fmt.Sprintf(" AND p.space_id=$%d", argIdx+1)
				args = append(args, tsquery, filter.SpaceID)
				argIdx += 2
			} else {
				args = append(args, tsquery)
				argIdx++
			}
			parts = append(parts, fmt.Sprintf(`
				SELECT 'page' AS type, p.id, p.title,
				       LEFT(p.content_text, 200) AS excerpt,
				       p.updated_at,
				       COALESCE(ts_rank(p.search_vector, plainto_tsquery('english', $%d)), 0) AS score
				FROM pages p
				WHERE p.deleted_at IS NULL AND p.status='published'
				  AND (p.search_vector @@ plainto_tsquery('english', $%d)
				       OR p.title ILIKE '%%' || $%d || '%%')
				  %s
			`, argIdx-1, argIdx-1, argIdx-1, cond))

		case entity.SearchResultProject:
			args = append(args, tsquery)
			argIdx++
			parts = append(parts, fmt.Sprintf(`
				SELECT 'project' AS type, pr.id, pr.name AS title,
				       LEFT(COALESCE(pr.description,''), 200) AS excerpt,
				       pr.updated_at,
				       0.0::float AS score
				FROM projects pr
				WHERE pr.deleted_at IS NULL AND pr.is_archived=FALSE
				  AND pr.name ILIKE '%%' || $%d || '%%'
			`, argIdx-1))

		case entity.SearchResultSpace:
			args = append(args, tsquery)
			argIdx++
			parts = append(parts, fmt.Sprintf(`
				SELECT 'space' AS type, sp.id, sp.name AS title,
				       '' AS excerpt,
				       sp.updated_at,
				       0.0::float AS score
				FROM spaces sp
				WHERE sp.deleted_at IS NULL AND sp.is_archived=FALSE
				  AND sp.name ILIKE '%%' || $%d || '%%'
			`, argIdx-1))
		}
	}

	if len(parts) == 0 {
		return nil, 0, nil
	}

	union := strings.Join(parts, " UNION ALL ")

	limit := filter.Limit
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	page := filter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	countSQL := fmt.Sprintf(`SELECT COUNT(*) FROM (%s) sub`, union)
	var total int
	if err := r.db.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("searchRepo.Search count: %w", err)
	}

	args = append(args, limit, offset)
	dataSQL := fmt.Sprintf(`
		SELECT type, id, title, excerpt, updated_at, score
		FROM (%s) sub
		ORDER BY score DESC, updated_at DESC
		LIMIT $%d OFFSET $%d
	`, union, argIdx, argIdx+1)

	rows, err := r.db.Query(ctx, dataSQL, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("searchRepo.Search query: %w", err)
	}
	defer rows.Close()

	var results []*entity.SearchResult
	for rows.Next() {
		res := &entity.SearchResult{}
		var updatedAt time.Time
		var score float64
		if err := rows.Scan(&res.Type, &res.ID, &res.Title, &res.Excerpt, &updatedAt, &score); err != nil {
			return nil, 0, err
		}
		res.UpdatedAt = updatedAt
		results = append(results, res)
	}
	return results, total, rows.Err()
}

func (r *searchRepo) Suggest(ctx context.Context, query string, limit int) ([]*entity.SearchSuggestion, error) {
	if limit <= 0 || limit > 20 {
		limit = 10
	}
	sql := `
		SELECT id, 'issue' AS type, title, word_similarity($1, title) AS score
		FROM issues WHERE title % $1 AND deleted_at IS NULL
		UNION ALL
		SELECT id, 'page' AS type, title, word_similarity($1, title) AS score
		FROM pages WHERE title % $1 AND deleted_at IS NULL
		UNION ALL
		SELECT id, 'project' AS type, name AS title, word_similarity($1, name) AS score
		FROM projects WHERE name % $1 AND deleted_at IS NULL
		UNION ALL
		SELECT id, 'space' AS type, name AS title, word_similarity($1, name) AS score
		FROM spaces WHERE name % $1 AND deleted_at IS NULL
		ORDER BY score DESC
		LIMIT $2
	`
	rows, err := r.db.Query(ctx, sql, query, limit)
	if err != nil {
		return nil, fmt.Errorf("searchRepo.Suggest: %w", err)
	}
	defer rows.Close()

	var suggestions []*entity.SearchSuggestion
	for rows.Next() {
		s := &entity.SearchSuggestion{}
		var score float64
		if err := rows.Scan(&s.ID, &s.Type, &s.Title, &score); err != nil {
			return nil, err
		}
		suggestions = append(suggestions, s)
	}
	return suggestions, rows.Err()
}
