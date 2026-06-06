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
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type projectTemplateRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewProjectTemplateRepo(p *pg.Postgres) repository.ProjectTemplateRepository {
	return &projectTemplateRepo{db: p.DB, builder: p.Builder}
}

func scanProjectTemplate(row pgx.Row) (*entity.ProjectTemplate, error) {
	t := &entity.ProjectTemplate{}
	var configJSON []byte
	if err := row.Scan(&t.ID, &t.Name, &t.Type, &t.Description, &t.IconURL, &configJSON, &t.IsSystem, &t.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.NotFound("project template not found")
		}
		return nil, fmt.Errorf("projectTemplateRepo.scan: %w", err)
	}
	if len(configJSON) > 0 {
		_ = json.Unmarshal(configJSON, &t.DefaultWorkflowConfig)
	}
	return t, nil
}

func (r *projectTemplateRepo) List(ctx context.Context) ([]*entity.ProjectTemplate, error) {
	sql, args, _ := r.builder.
		Select("id, name, type, description, icon_url, default_workflow_config, is_system, created_at").
		From("project_templates").OrderBy("is_system DESC, name ASC").ToSql()
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("projectTemplateRepo.List: %w", err)
	}
	defer rows.Close()
	var list []*entity.ProjectTemplate
	for rows.Next() {
		t, err := scanProjectTemplate(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, rows.Err()
}

func (r *projectTemplateRepo) GetByID(ctx context.Context, id string) (*entity.ProjectTemplate, error) {
	sql, args, _ := r.builder.
		Select("id, name, type, description, icon_url, default_workflow_config, is_system, created_at").
		From("project_templates").Where(sq.Eq{"id": id}).ToSql()
	return scanProjectTemplate(r.db.QueryRow(ctx, sql, args...))
}
