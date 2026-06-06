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

type spaceCategoryRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewSpaceCategoryRepo(p *pg.Postgres) repository.SpaceCategoryRepository {
	return &spaceCategoryRepo{db: p.DB, builder: p.Builder}
}

func scanSpaceCategory(row pgx.Row) (*entity.SpaceCategory, error) {
	c := &entity.SpaceCategory{}
	if err := row.Scan(&c.ID, &c.Name, &c.Color, &c.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.NotFound("space category not found")
		}
		return nil, fmt.Errorf("spaceCategoryRepo.scan: %w", err)
	}
	return c, nil
}

func (r *spaceCategoryRepo) Create(ctx context.Context, c *entity.SpaceCategory) error {
	sql, args, err := r.builder.
		Insert("space_categories").
		Columns("id", "name", "color").
		Values(c.ID, c.Name, c.Color).
		ToSql()
	if err != nil {
		return fmt.Errorf("spaceCategoryRepo.Create build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *spaceCategoryRepo) GetByID(ctx context.Context, id string) (*entity.SpaceCategory, error) {
	sql, args, _ := r.builder.Select("id, name, color, created_at").From("space_categories").Where(sq.Eq{"id": id}).ToSql()
	return scanSpaceCategory(r.db.QueryRow(ctx, sql, args...))
}

func (r *spaceCategoryRepo) List(ctx context.Context) ([]*entity.SpaceCategory, error) {
	sql, args, _ := r.builder.Select("id, name, color, created_at").From("space_categories").OrderBy("name ASC").ToSql()
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("spaceCategoryRepo.List: %w", err)
	}
	defer rows.Close()
	var list []*entity.SpaceCategory
	for rows.Next() {
		c, err := scanSpaceCategory(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func (r *spaceCategoryRepo) Update(ctx context.Context, c *entity.SpaceCategory) error {
	sql, args, err := r.builder.Update("space_categories").
		Set("name", c.Name).Set("color", c.Color).
		Where(sq.Eq{"id": c.ID}).ToSql()
	if err != nil {
		return fmt.Errorf("spaceCategoryRepo.Update build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *spaceCategoryRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM space_categories WHERE id=$1`, id)
	return err
}
