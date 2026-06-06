package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type blueprintRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewBlueprintRepo(p *pg.Postgres) repository.BlueprintRepository {
	return &blueprintRepo{db: p.DB, builder: p.Builder}
}

func scanBlueprint(row pgx.Row) (*entity.Blueprint, error) {
	b := &entity.Blueprint{}
	var schemaJSON []byte
	if err := row.Scan(&b.ID, &b.Name, &b.Description, &b.IconURL, &b.Category,
		&b.TemplateBody, &schemaJSON, &b.IsSystem, &b.CreatedAt, &b.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.NotFound("blueprint not found")
		}
		return nil, fmt.Errorf("blueprintRepo.scan: %w", err)
	}
	if len(schemaJSON) > 0 {
		_ = json.Unmarshal(schemaJSON, &b.Schema)
	}
	return b, nil
}

func (r *blueprintRepo) Create(ctx context.Context, b *entity.Blueprint) error {
	schemaJSON, _ := json.Marshal(b.Schema)
	sql, args, err := r.builder.
		Insert("blueprints").
		Columns("id", "name", "description", "icon_url", "category", "template_body", "schema", "is_system").
		Values(b.ID, b.Name, b.Description, b.IconURL, b.Category, b.TemplateBody, schemaJSON, b.IsSystem).
		ToSql()
	if err != nil {
		return fmt.Errorf("blueprintRepo.Create build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *blueprintRepo) GetByID(ctx context.Context, id string) (*entity.Blueprint, error) {
	sql, args, _ := r.builder.
		Select("id, name, description, icon_url, category, template_body, schema, is_system, created_at, updated_at").
		From("blueprints").Where(sq.Eq{"id": id}).ToSql()
	return scanBlueprint(r.db.QueryRow(ctx, sql, args...))
}

func (r *blueprintRepo) List(ctx context.Context) ([]*entity.Blueprint, error) {
	sql, args, _ := r.builder.
		Select("id, name, description, icon_url, category, template_body, schema, is_system, created_at, updated_at").
		From("blueprints").OrderBy("is_system DESC, name ASC").ToSql()
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("blueprintRepo.List: %w", err)
	}
	defer rows.Close()
	var list []*entity.Blueprint
	for rows.Next() {
		b, err := scanBlueprint(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, b)
	}
	return list, rows.Err()
}

func (r *blueprintRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM blueprints WHERE id=$1 AND is_system=false`, id)
	return err
}
