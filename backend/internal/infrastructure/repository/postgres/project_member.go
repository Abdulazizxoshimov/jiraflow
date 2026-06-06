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

type projectMemberRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewProjectMemberRepo(p *pg.Postgres) repository.ProjectMemberRepository {
	return &projectMemberRepo{db: p.DB, builder: p.Builder}
}

func (r *projectMemberRepo) Add(ctx context.Context, m *entity.ProjectMember) error {
	sql, args, err := r.builder.
		Insert("project_members").
		Columns("project_id", "user_id", "role", "created_at").
		Values(m.ProjectID, m.UserID, m.Role, m.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("projectMemberRepo.Add: %w", err)
	}
	if _, err = r.db.Exec(ctx, sql, args...); err != nil {
		if isUniqueViolation(err) {
			return apperr.Conflict("user is already a member of this project")
		}
		return fmt.Errorf("projectMemberRepo.Add: %w", err)
	}
	return nil
}

func (r *projectMemberRepo) GetMember(ctx context.Context, projectID, userID string) (*entity.ProjectMember, error) {
	const q = `SELECT pm.project_id, pm.user_id, pm.role, pm.created_at,
		u.id, u.full_name, u.avatar_url, u.email, u.color
		FROM project_members pm
		JOIN users u ON u.id = pm.user_id
		WHERE pm.project_id = $1 AND pm.user_id = $2`

	m := &entity.ProjectMember{User: &entity.UserShort{}}
	err := r.db.QueryRow(ctx, q, projectID, userID).Scan(
		&m.ProjectID, &m.UserID, &m.Role, &m.CreatedAt,
		&m.User.ID, &m.User.FullName, &m.User.AvatarURL, &m.User.Email, &m.User.Color,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("project member")
	}
	return m, err
}

func (r *projectMemberRepo) ListByProject(ctx context.Context, projectID string, filter *entity.Filter) ([]*entity.ProjectMember, int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM project_members WHERE project_id = $1`, projectID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("projectMemberRepo.ListByProject count: %w", err)
	}

	const q = `SELECT pm.project_id, pm.user_id, pm.role, pm.created_at,
		u.id, u.full_name, u.avatar_url, u.email, u.color
		FROM project_members pm
		JOIN users u ON u.id = pm.user_id
		WHERE pm.project_id = $1
		ORDER BY pm.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, q, projectID, filter.GetLimit(), filter.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("projectMemberRepo.ListByProject query: %w", err)
	}
	defer rows.Close()

	var members []*entity.ProjectMember
	for rows.Next() {
		m := &entity.ProjectMember{User: &entity.UserShort{}}
		if err := rows.Scan(
			&m.ProjectID, &m.UserID, &m.Role, &m.CreatedAt,
			&m.User.ID, &m.User.FullName, &m.User.AvatarURL, &m.User.Email, &m.User.Color,
		); err != nil {
			return nil, 0, err
		}
		members = append(members, m)
	}
	return members, total, rows.Err()
}

func (r *projectMemberRepo) ListByUser(ctx context.Context, userID string) ([]*entity.ProjectMember, error) {
	sql, args, err := r.builder.
		Select("project_id", "user_id", "role", "created_at").
		From("project_members").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("projectMemberRepo.ListByUser: %w", err)
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("projectMemberRepo.ListByUser query: %w", err)
	}
	defer rows.Close()

	var members []*entity.ProjectMember
	for rows.Next() {
		m := &entity.ProjectMember{}
		if err := rows.Scan(&m.ProjectID, &m.UserID, &m.Role, &m.CreatedAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

func (r *projectMemberRepo) UpdateRole(ctx context.Context, projectID, userID, role string) error {
	sql, args, err := r.builder.
		Update("project_members").
		Set("role", role).
		Where(sq.And{sq.Eq{"project_id": projectID}, sq.Eq{"user_id": userID}}).
		ToSql()
	if err != nil {
		return fmt.Errorf("projectMemberRepo.UpdateRole: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *projectMemberRepo) Remove(ctx context.Context, projectID, userID string) error {
	sql, args, err := r.builder.
		Delete("project_members").
		Where(sq.And{sq.Eq{"project_id": projectID}, sq.Eq{"user_id": userID}}).
		ToSql()
	if err != nil {
		return fmt.Errorf("projectMemberRepo.Remove: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *projectMemberRepo) IsMember(ctx context.Context, projectID, userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM project_members WHERE project_id=$1 AND user_id=$2)`,
		projectID, userID,
	).Scan(&exists)
	return exists, err
}
