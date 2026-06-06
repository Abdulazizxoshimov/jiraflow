package postgres

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type fieldConfigRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewFieldConfigurationRepo(p *pg.Postgres) repository.FieldConfigurationRepository {
	return &fieldConfigRepo{db: p.DB, builder: p.Builder}
}

func (r *fieldConfigRepo) Create(ctx context.Context, c *entity.FieldConfiguration) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("fieldConfigRepo.Create begin: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO field_configurations (id, name, project_id) VALUES ($1, $2, $3)`,
		c.ID, c.Name, c.ProjectID)
	if err != nil {
		return fmt.Errorf("fieldConfigRepo.Create config: %w", err)
	}

	for _, item := range c.Items {
		if item.ID == "" {
			item.ID = uuid.NewString()
		}
		item.ConfigID = c.ID
		_, err = tx.Exec(ctx,
			`INSERT INTO field_config_items (id, config_id, field_name, is_required, is_hidden, description) VALUES ($1,$2,$3,$4,$5,$6)`,
			item.ID, item.ConfigID, item.FieldName, item.IsRequired, item.IsHidden, item.Description)
		if err != nil {
			return fmt.Errorf("fieldConfigRepo.Create item: %w", err)
		}
	}
	return tx.Commit(ctx)
}

func (r *fieldConfigRepo) GetByID(ctx context.Context, id string) (*entity.FieldConfiguration, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, name, project_id, created_at FROM field_configurations WHERE id=$1`, id)
	c := &entity.FieldConfiguration{}
	if err := row.Scan(&c.ID, &c.Name, &c.ProjectID, &c.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.NotFound("field configuration not found")
		}
		return nil, fmt.Errorf("fieldConfigRepo.GetByID: %w", err)
	}
	rows, err := r.db.Query(ctx,
		`SELECT id, config_id, field_name, is_required, is_hidden, description FROM field_config_items WHERE config_id=$1 ORDER BY field_name`, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			item := &entity.FieldConfigItem{}
			if rows.Scan(&item.ID, &item.ConfigID, &item.FieldName, &item.IsRequired, &item.IsHidden, &item.Description) == nil {
				c.Items = append(c.Items, item)
			}
		}
	}
	return c, nil
}

func (r *fieldConfigRepo) List(ctx context.Context, projectID *string) ([]*entity.FieldConfiguration, error) {
	q := r.builder.Select("id, name, project_id, created_at").From("field_configurations")
	if projectID != nil {
		q = q.Where(sq.Eq{"project_id": *projectID})
	}
	sql, args, _ := q.OrderBy("created_at DESC").ToSql()
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("fieldConfigRepo.List: %w", err)
	}
	defer rows.Close()
	var list []*entity.FieldConfiguration
	for rows.Next() {
		c := &entity.FieldConfiguration{}
		if err := rows.Scan(&c.ID, &c.Name, &c.ProjectID, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("fieldConfigRepo.List scan: %w", err)
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func (r *fieldConfigRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM field_configurations WHERE id=$1`, id)
	return err
}
