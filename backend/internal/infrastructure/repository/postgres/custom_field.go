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

type customFieldRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewCustomFieldRepo(p *pg.Postgres) repository.CustomFieldRepository {
	return &customFieldRepo{db: p.DB, builder: p.Builder}
}

func scanCustomField(row pgx.Row) (*entity.CustomField, error) {
	cf := &entity.CustomField{}
	var optionsJSON []byte
	err := row.Scan(
		&cf.ID, &cf.ProjectID, &cf.Name, &cf.FieldKey, &cf.FieldType,
		&cf.IsRequired, &optionsJSON, &cf.Position, &cf.CreatedAt, &cf.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if len(optionsJSON) > 0 {
		_ = json.Unmarshal(optionsJSON, &cf.Options)
	}
	return cf, nil
}

const cfCols = "id, project_id, name, field_key, field_type, is_required, options, position, created_at, updated_at"

func (r *customFieldRepo) Create(ctx context.Context, cf *entity.CustomField) error {
	optJSON, err := json.Marshal(cf.Options)
	if err != nil {
		return fmt.Errorf("customFieldRepo.Create marshal options: %w", err)
	}
	sql, args, err := r.builder.
		Insert("custom_fields").
		Columns("id", "project_id", "name", "field_key", "field_type", "is_required", "options", "position", "created_at", "updated_at").
		Values(cf.ID, cf.ProjectID, cf.Name, cf.FieldKey, cf.FieldType, cf.IsRequired, optJSON, cf.Position, cf.CreatedAt, cf.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("customFieldRepo.Create: %w", err)
	}
	if _, err = r.db.Exec(ctx, sql, args...); err != nil {
		if isUniqueViolation(err) {
			return apperr.Conflict("field key already exists in this project")
		}
		return fmt.Errorf("customFieldRepo.Create: %w", err)
	}
	return nil
}

func (r *customFieldRepo) GetByID(ctx context.Context, id string) (*entity.CustomField, error) {
	sql, args, err := r.builder.
		Select(cfCols).From("custom_fields").
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("customFieldRepo.GetByID: %w", err)
	}
	cf, err := scanCustomField(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("custom field")
	}
	return cf, err
}

func (r *customFieldRepo) GetByKey(ctx context.Context, projectID, fieldKey string) (*entity.CustomField, error) {
	sql, args, err := r.builder.
		Select(cfCols).From("custom_fields").
		Where(sq.And{sq.Eq{"project_id": projectID}, sq.Eq{"field_key": fieldKey}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("customFieldRepo.GetByKey: %w", err)
	}
	cf, err := scanCustomField(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("custom field")
	}
	return cf, err
}

func (r *customFieldRepo) ListByProject(ctx context.Context, projectID string) ([]*entity.CustomField, error) {
	sql, args, err := r.builder.
		Select(cfCols).From("custom_fields").
		Where(sq.And{sq.Eq{"project_id": projectID}, sq.Eq{"deleted_at": nil}}).
		OrderBy("position ASC").ToSql()
	if err != nil {
		return nil, fmt.Errorf("customFieldRepo.ListByProject: %w", err)
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("customFieldRepo.ListByProject query: %w", err)
	}
	defer rows.Close()

	var fields []*entity.CustomField
	for rows.Next() {
		cf, err := scanCustomField(rows)
		if err != nil {
			return nil, err
		}
		fields = append(fields, cf)
	}
	return fields, rows.Err()
}

func (r *customFieldRepo) Update(ctx context.Context, cf *entity.CustomField) error {
	optJSON, err := json.Marshal(cf.Options)
	if err != nil {
		return fmt.Errorf("customFieldRepo.Update marshal options: %w", err)
	}
	sql, args, err := r.builder.
		Update("custom_fields").
		Set("name", cf.Name).Set("is_required", cf.IsRequired).
		Set("options", optJSON).Set("position", cf.Position).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": cf.ID}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return fmt.Errorf("customFieldRepo.Update: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *customFieldRepo) Delete(ctx context.Context, id string) error {
	sql, args, err := r.builder.
		Update("custom_fields").Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return fmt.Errorf("customFieldRepo.Delete: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *customFieldRepo) ReorderFields(ctx context.Context, projectID string, positions map[string]int) error {
	for id, pos := range positions {
		sql, args, err := r.builder.
			Update("custom_fields").
			Set("position", pos).Set("updated_at", sq.Expr("NOW()")).
			Where(sq.And{sq.Eq{"id": id}, sq.Eq{"project_id": projectID}}).ToSql()
		if err != nil {
			return fmt.Errorf("customFieldRepo.ReorderFields: %w", err)
		}
		if _, err = r.db.Exec(ctx, sql, args...); err != nil {
			return fmt.Errorf("customFieldRepo.ReorderFields exec: %w", err)
		}
	}
	return nil
}
