package postgres

import (
	"context"
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

type issueTypeRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewIssueTypeRepo(p *pg.Postgres) repository.IssueTypeRepository {
	return &issueTypeRepo{db: p.DB, builder: p.Builder}
}

func scanIssueType(row pgx.Row) (*entity.IssueType, error) {
	t := &entity.IssueType{}
	if err := row.Scan(&t.ID, &t.Name, &t.Description, &t.IconURL, &t.Color, &t.IsSubtask, &t.IsSystem, &t.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.NotFound("issue type not found")
		}
		return nil, fmt.Errorf("issueTypeRepo.scan: %w", err)
	}
	return t, nil
}

func (r *issueTypeRepo) CreateType(ctx context.Context, t *entity.IssueType) error {
	sql, args, err := r.builder.Insert("issue_types").
		Columns("id", "name", "description", "icon_url", "color", "is_subtask", "is_system").
		Values(t.ID, t.Name, t.Description, t.IconURL, t.Color, t.IsSubtask, t.IsSystem).ToSql()
	if err != nil {
		return fmt.Errorf("issueTypeRepo.CreateType build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *issueTypeRepo) ListTypes(ctx context.Context) ([]*entity.IssueType, error) {
	sql, args, _ := r.builder.Select("id, name, description, icon_url, color, is_subtask, is_system, created_at").
		From("issue_types").OrderBy("is_system DESC, name ASC").ToSql()
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("issueTypeRepo.ListTypes: %w", err)
	}
	defer rows.Close()
	var list []*entity.IssueType
	for rows.Next() {
		t, err := scanIssueType(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, rows.Err()
}

func (r *issueTypeRepo) GetTypeByID(ctx context.Context, id string) (*entity.IssueType, error) {
	sql, args, _ := r.builder.Select("id, name, description, icon_url, color, is_subtask, is_system, created_at").
		From("issue_types").Where(sq.Eq{"id": id}).ToSql()
	return scanIssueType(r.db.QueryRow(ctx, sql, args...))
}

func (r *issueTypeRepo) DeleteType(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM issue_types WHERE id=$1 AND is_system=false`, id)
	return err
}

func (r *issueTypeRepo) CreateScheme(ctx context.Context, s *entity.IssueTypeScheme, issueTypeIDs []string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("issueTypeRepo.CreateScheme begin: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `INSERT INTO issue_type_schemes (id, name, project_id) VALUES ($1, $2, $3)`,
		s.ID, s.Name, s.ProjectID)
	if err != nil {
		return fmt.Errorf("issueTypeRepo.CreateScheme insert scheme: %w", err)
	}

	for i, typeID := range issueTypeIDs {
		_, err = tx.Exec(ctx, `INSERT INTO issue_type_scheme_members (scheme_id, issue_type_id, display_order) VALUES ($1, $2, $3)`,
			s.ID, typeID, i)
		if err != nil {
			return fmt.Errorf("issueTypeRepo.CreateScheme insert member: %w", err)
		}
	}
	return tx.Commit(ctx)
}

func (r *issueTypeRepo) GetSchemeByID(ctx context.Context, id string) (*entity.IssueTypeScheme, error) {
	return r.loadScheme(ctx, `WHERE s.id=$1`, id)
}

func (r *issueTypeRepo) GetSchemeByProject(ctx context.Context, projectID string) (*entity.IssueTypeScheme, error) {
	return r.loadScheme(ctx, `WHERE s.project_id=$1`, projectID)
}

func (r *issueTypeRepo) loadScheme(ctx context.Context, whereClause string, arg any) (*entity.IssueTypeScheme, error) {
	row := r.db.QueryRow(ctx,
		`SELECT s.id, s.name, s.project_id, s.created_at FROM issue_type_schemes s `+whereClause, arg)
	s := &entity.IssueTypeScheme{}
	if err := row.Scan(&s.ID, &s.Name, &s.ProjectID, &s.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.NotFound("issue type scheme not found")
		}
		return nil, fmt.Errorf("issueTypeRepo.loadScheme: %w", err)
	}

	rows, err := r.db.Query(ctx, `
		SELECT t.id, t.name, t.description, t.icon_url, t.color, t.is_subtask, t.is_system, t.created_at
		FROM issue_type_scheme_members m
		JOIN issue_types t ON t.id = m.issue_type_id
		WHERE m.scheme_id=$1 ORDER BY m.display_order ASC`, s.ID)
	if err != nil {
		return nil, fmt.Errorf("issueTypeRepo.loadScheme members: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		t, err := scanIssueType(rows)
		if err != nil {
			return nil, err
		}
		s.IssueTypes = append(s.IssueTypes, t)
	}
	return s, rows.Err()
}

func (r *issueTypeRepo) ListSchemes(ctx context.Context) ([]*entity.IssueTypeScheme, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, project_id, created_at FROM issue_type_schemes ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("issueTypeRepo.ListSchemes: %w", err)
	}
	defer rows.Close()
	var list []*entity.IssueTypeScheme
	for rows.Next() {
		s := &entity.IssueTypeScheme{}
		if err := rows.Scan(&s.ID, &s.Name, &s.ProjectID, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("issueTypeRepo.ListSchemes scan: %w", err)
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

func (r *issueTypeRepo) DeleteScheme(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM issue_type_schemes WHERE id=$1`, id)
	return err
}
