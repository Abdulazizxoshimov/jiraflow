package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type savedFilterRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewSavedFilterRepo(p *pg.Postgres) repository.SavedFilterRepository {
	return &savedFilterRepo{db: p.DB, builder: p.Builder}
}

func scanSavedFilter(row pgx.Row) (*entity.SavedFilter, error) {
	sf := &entity.SavedFilter{}
	var filtersJSON []byte
	err := row.Scan(
		&sf.ID, &sf.UserID, &sf.Name, &sf.Description,
		&sf.FilterType, &filtersJSON, &sf.IsShared,
		&sf.CreatedAt, &sf.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if len(filtersJSON) > 0 {
		_ = json.Unmarshal(filtersJSON, &sf.Filters)
	}
	return sf, nil
}

const savedFilterCols = "id, user_id, name, description, filter_type, filters, is_shared, created_at, updated_at"

func (r *savedFilterRepo) Create(ctx context.Context, userID string, req *entity.CreateSavedFilterReq) (*entity.SavedFilter, error) {
	filtersJSON, err := json.Marshal(req.Filters)
	if err != nil {
		return nil, fmt.Errorf("savedFilterRepo.Create marshal: %w", err)
	}
	id := uuid.NewString()
	_, err = r.db.Exec(ctx, `
		INSERT INTO saved_filters(id, user_id, name, description, filter_type, filters, is_shared)
		VALUES($1,$2,$3,$4,$5,$6,$7)
	`, id, userID, req.Name, req.Description, req.FilterType, filtersJSON, req.IsShared)
	if err != nil {
		return nil, fmt.Errorf("savedFilterRepo.Create: %w", err)
	}
	return r.GetByID(ctx, id)
}

func (r *savedFilterRepo) GetByID(ctx context.Context, id string) (*entity.SavedFilter, error) {
	sql, args, _ := r.builder.Select(savedFilterCols).From("saved_filters").Where(sq.Eq{"id": id}).ToSql()
	sf, err := scanSavedFilter(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("saved filter")
	}
	return sf, err
}

func (r *savedFilterRepo) List(ctx context.Context, userID, filterType string) ([]*entity.SavedFilter, error) {
	where := sq.Or{sq.Eq{"user_id": userID}, sq.Eq{"is_shared": true}}
	q := r.builder.Select(savedFilterCols).From("saved_filters").Where(where)
	if filterType != "" {
		q = q.Where(sq.Eq{"filter_type": filterType})
	}
	sql, args, _ := q.OrderBy("created_at DESC").ToSql()

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("savedFilterRepo.List: %w", err)
	}
	defer rows.Close()

	var filters []*entity.SavedFilter
	for rows.Next() {
		sf, err := scanSavedFilter(rows)
		if err != nil {
			return nil, err
		}
		filters = append(filters, sf)
	}
	return filters, rows.Err()
}

func (r *savedFilterRepo) Update(ctx context.Context, id string, req *entity.UpdateSavedFilterReq) (*entity.SavedFilter, error) {
	q := r.builder.Update("saved_filters").Where(sq.Eq{"id": id})
	if req.Name != nil {
		q = q.Set("name", *req.Name)
	}
	if req.Description != nil {
		q = q.Set("description", *req.Description)
	}
	if req.Filters != nil {
		filtersJSON, _ := json.Marshal(req.Filters)
		q = q.Set("filters", filtersJSON)
	}
	if req.IsShared != nil {
		q = q.Set("is_shared", *req.IsShared)
	}
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("savedFilterRepo.Update build: %w", err)
	}
	if _, err := r.db.Exec(ctx, sql, args...); err != nil {
		return nil, fmt.Errorf("savedFilterRepo.Update: %w", err)
	}
	return r.GetByID(ctx, id)
}

func (r *savedFilterRepo) Delete(ctx context.Context, id, userID string) error {
	sql, args, err := r.builder.Delete("saved_filters").
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"user_id": userID}}).ToSql()
	if err != nil {
		return fmt.Errorf("savedFilterRepo.Delete build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}
