package postgres

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type spaceExportRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewSpaceExportRepo(p *pg.Postgres) repository.SpaceExportRepository {
	return &spaceExportRepo{db: p.DB, builder: p.Builder}
}

const spaceExportCols = "id, space_id, requested_by, status, file_url, error_msg, created_at, updated_at"

func scanSpaceExport(row pgx.Row) (*entity.SpaceExport, error) {
	e := &entity.SpaceExport{}
	if err := row.Scan(&e.ID, &e.SpaceID, &e.RequestedBy, &e.Status, &e.FileURL, &e.ErrorMsg, &e.CreatedAt, &e.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.NotFound("space export not found")
		}
		return nil, fmt.Errorf("spaceExportRepo.scan: %w", err)
	}
	return e, nil
}

func (r *spaceExportRepo) Create(ctx context.Context, e *entity.SpaceExport) error {
	sql, args, err := r.builder.
		Insert("space_exports").
		Columns("id", "space_id", "requested_by", "status").
		Values(e.ID, e.SpaceID, e.RequestedBy, e.Status).
		ToSql()
	if err != nil {
		return fmt.Errorf("spaceExportRepo.Create build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *spaceExportRepo) GetByID(ctx context.Context, id string) (*entity.SpaceExport, error) {
	sql, args, _ := r.builder.Select(spaceExportCols).From("space_exports").Where(sq.Eq{"id": id}).ToSql()
	return scanSpaceExport(r.db.QueryRow(ctx, sql, args...))
}

func (r *spaceExportRepo) UpdateStatus(ctx context.Context, id, status string, fileURL, errorMsg *string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE space_exports SET status=$1, file_url=$2, error_msg=$3, updated_at=NOW() WHERE id=$4`,
		status, fileURL, errorMsg, id)
	return err
}

func (r *spaceExportRepo) ListBySpace(ctx context.Context, spaceID string) ([]*entity.SpaceExport, error) {
	sql, args, _ := r.builder.Select(spaceExportCols).From("space_exports").
		Where(sq.Eq{"space_id": spaceID}).OrderBy("created_at DESC").Limit(20).ToSql()
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("spaceExportRepo.ListBySpace: %w", err)
	}
	defer rows.Close()
	var list []*entity.SpaceExport
	for rows.Next() {
		e, err := scanSpaceExport(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}
