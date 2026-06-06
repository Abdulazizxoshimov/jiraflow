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

type permissionSchemeRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewPermissionSchemeRepo(p *pg.Postgres) repository.PermissionSchemeRepository {
	return &permissionSchemeRepo{db: p.DB, builder: p.Builder}
}

func (r *permissionSchemeRepo) Create(ctx context.Context, s *entity.PermissionScheme) error {
	sql, args, err := r.builder.
		Insert("permission_schemes").
		Columns("id", "name", "description", "created_by", "created_at", "updated_at").
		Values(s.ID, s.Name, s.Description, s.CreatedBy, s.CreatedAt, s.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("permissionSchemeRepo.Create build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *permissionSchemeRepo) GetByID(ctx context.Context, id string) (*entity.PermissionScheme, error) {
	sql, args, err := r.builder.
		Select("id", "name", "description", "created_by", "created_at", "updated_at").
		From("permission_schemes").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("permissionSchemeRepo.GetByID build: %w", err)
	}
	var s entity.PermissionScheme
	err = r.db.QueryRow(ctx, sql, args...).Scan(&s.ID, &s.Name, &s.Description, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("permission scheme not found")
	}
	if err != nil {
		return nil, fmt.Errorf("permissionSchemeRepo.GetByID scan: %w", err)
	}
	grants, _ := r.ListGrants(ctx, id)
	s.Grants = grants
	return &s, nil
}

func (r *permissionSchemeRepo) List(ctx context.Context) ([]*entity.PermissionScheme, error) {
	sql, args, err := r.builder.
		Select("id", "name", "description", "created_by", "created_at", "updated_at").
		From("permission_schemes").
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("permissionSchemeRepo.List build: %w", err)
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("permissionSchemeRepo.List query: %w", err)
	}
	defer rows.Close()
	var list []*entity.PermissionScheme
	for rows.Next() {
		var s entity.PermissionScheme
		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("permissionSchemeRepo.List scan: %w", err)
		}
		list = append(list, &s)
	}
	return list, nil
}

func (r *permissionSchemeRepo) Update(ctx context.Context, s *entity.PermissionScheme) error {
	sql, args, err := r.builder.
		Update("permission_schemes").
		Set("name", s.Name).
		Set("description", s.Description).
		Set("updated_at", s.UpdatedAt).
		Where(sq.Eq{"id": s.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("permissionSchemeRepo.Update build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *permissionSchemeRepo) Delete(ctx context.Context, id string) error {
	sql, args, err := r.builder.Delete("permission_schemes").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("permissionSchemeRepo.Delete build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *permissionSchemeRepo) AddGrant(ctx context.Context, g *entity.PermissionSchemeGrant) error {
	sql, args, err := r.builder.
		Insert("permission_scheme_grants").
		Columns("id", "scheme_id", "permission", "holder_type", "holder_id", "created_at").
		Values(g.ID, g.SchemeID, g.Permission, g.HolderType, g.HolderID, g.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("permissionSchemeRepo.AddGrant build: %w", err)
	}
	if _, err = r.db.Exec(ctx, sql, args...); err != nil {
		if isUniqueViolation(err) {
			return apperr.Conflict("grant already exists")
		}
		return fmt.Errorf("permissionSchemeRepo.AddGrant exec: %w", err)
	}
	return nil
}

func (r *permissionSchemeRepo) RemoveGrant(ctx context.Context, grantID string) error {
	sql, args, err := r.builder.Delete("permission_scheme_grants").Where(sq.Eq{"id": grantID}).ToSql()
	if err != nil {
		return fmt.Errorf("permissionSchemeRepo.RemoveGrant build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *permissionSchemeRepo) ListGrants(ctx context.Context, schemeID string) ([]*entity.PermissionSchemeGrant, error) {
	sql, args, err := r.builder.
		Select("id", "scheme_id", "permission", "holder_type", "holder_id", "created_at").
		From("permission_scheme_grants").
		Where(sq.Eq{"scheme_id": schemeID}).
		OrderBy("created_at ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("permissionSchemeRepo.ListGrants build: %w", err)
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("permissionSchemeRepo.ListGrants query: %w", err)
	}
	defer rows.Close()
	var list []*entity.PermissionSchemeGrant
	for rows.Next() {
		var g entity.PermissionSchemeGrant
		if err := rows.Scan(&g.ID, &g.SchemeID, &g.Permission, &g.HolderType, &g.HolderID, &g.CreatedAt); err != nil {
			return nil, fmt.Errorf("permissionSchemeRepo.ListGrants scan: %w", err)
		}
		list = append(list, &g)
	}
	return list, nil
}

func (r *permissionSchemeRepo) AssignToProject(ctx context.Context, projectID, schemeID string) error {
	sql, args, err := r.builder.
		Update("projects").
		Set("permission_scheme_id", schemeID).
		Where(sq.Eq{"id": projectID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("permissionSchemeRepo.AssignToProject build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *permissionSchemeRepo) GetByProject(ctx context.Context, projectID string) (*entity.PermissionScheme, error) {
	sql, args, err := r.builder.
		Select("ps.id", "ps.name", "ps.description", "ps.created_by", "ps.created_at", "ps.updated_at").
		From("permission_schemes ps").
		Join("projects p ON p.permission_scheme_id = ps.id").
		Where(sq.Eq{"p.id": projectID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("permissionSchemeRepo.GetByProject build: %w", err)
	}
	var s entity.PermissionScheme
	err = r.db.QueryRow(ctx, sql, args...).Scan(&s.ID, &s.Name, &s.Description, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("no permission scheme assigned to project")
	}
	if err != nil {
		return nil, fmt.Errorf("permissionSchemeRepo.GetByProject scan: %w", err)
	}
	grants, _ := r.ListGrants(ctx, s.ID)
	s.Grants = grants
	return &s, nil
}
