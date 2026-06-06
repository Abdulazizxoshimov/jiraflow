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

type labelRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewLabelRepo(p *pg.Postgres) repository.LabelRepository {
	return &labelRepo{db: p.DB, builder: p.Builder}
}

func scanLabel(row pgx.Row) (*entity.Label, error) {
	l := &entity.Label{}
	err := row.Scan(&l.ID, &l.ProjectID, &l.Name, &l.Color, &l.CreatedAt)
	return l, err
}

func (r *labelRepo) Create(ctx context.Context, l *entity.Label) error {
	sql, args, err := r.builder.
		Insert("labels").
		Columns("id", "project_id", "name", "color", "created_at").
		Values(l.ID, l.ProjectID, l.Name, l.Color, l.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("labelRepo.Create: %w", err)
	}
	if _, err = r.db.Exec(ctx, sql, args...); err != nil {
		if isUniqueViolation(err) {
			return apperr.Conflict("label name already exists in this project")
		}
		return fmt.Errorf("labelRepo.Create: %w", err)
	}
	return nil
}

func (r *labelRepo) GetByID(ctx context.Context, id string) (*entity.Label, error) {
	sql, args, err := r.builder.
		Select("id", "project_id", "name", "color", "created_at").
		From("labels").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("labelRepo.GetByID: %w", err)
	}
	l, err := scanLabel(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("label")
	}
	return l, err
}

func (r *labelRepo) ListByProject(ctx context.Context, projectID string) ([]*entity.Label, error) {
	sql, args, err := r.builder.
		Select("id", "project_id", "name", "color", "created_at").
		From("labels").Where(sq.Eq{"project_id": projectID}).
		OrderBy("name ASC").ToSql()
	if err != nil {
		return nil, fmt.Errorf("labelRepo.ListByProject: %w", err)
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("labelRepo.ListByProject query: %w", err)
	}
	defer rows.Close()

	var labels []*entity.Label
	for rows.Next() {
		l, err := scanLabel(rows)
		if err != nil {
			return nil, err
		}
		labels = append(labels, l)
	}
	return labels, rows.Err()
}

func (r *labelRepo) Update(ctx context.Context, l *entity.Label) error {
	sql, args, err := r.builder.
		Update("labels").
		Set("name", l.Name).Set("color", l.Color).
		Where(sq.Eq{"id": l.ID}).ToSql()
	if err != nil {
		return fmt.Errorf("labelRepo.Update: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *labelRepo) Delete(ctx context.Context, id string) error {
	sql, args, err := r.builder.Delete("labels").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("labelRepo.Delete: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *labelRepo) GetByIDs(ctx context.Context, ids []string) ([]*entity.Label, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	sql, args, err := r.builder.
		Select("id", "project_id", "name", "color", "created_at").
		From("labels").Where(sq.Eq{"id": ids}).OrderBy("name ASC").ToSql()
	if err != nil {
		return nil, fmt.Errorf("labelRepo.GetByIDs: %w", err)
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("labelRepo.GetByIDs query: %w", err)
	}
	defer rows.Close()

	var labels []*entity.Label
	for rows.Next() {
		l, err := scanLabel(rows)
		if err != nil {
			return nil, err
		}
		labels = append(labels, l)
	}
	return labels, rows.Err()
}
