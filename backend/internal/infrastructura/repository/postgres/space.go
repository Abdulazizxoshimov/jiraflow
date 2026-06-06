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

type spaceRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewSpaceRepo(p *pg.Postgres) repository.SpaceRepository {
	return &spaceRepo{db: p.DB, builder: p.Builder}
}

const spaceCols = "id, key, name, description, icon_url, type, lead_id, project_id, is_archived, created_at, updated_at, deleted_at"

func scanSpace(row pgx.Row) (*entity.Space, error) {
	s := &entity.Space{}
	err := row.Scan(
		&s.ID, &s.Key, &s.Name, &s.Description, &s.IconURL,
		&s.Type, &s.LeadID, &s.ProjectID, &s.IsArchived,
		&s.CreatedAt, &s.UpdatedAt, &s.DeletedAt,
	)
	return s, err
}

func (r *spaceRepo) Create(ctx context.Context, s *entity.Space) error {
	sql, args, err := r.builder.
		Insert("spaces").
		Columns("id", "key", "name", "description", "icon_url", "type", "lead_id", "project_id", "created_at", "updated_at").
		Values(s.ID, s.Key, s.Name, s.Description, s.IconURL, s.Type, s.LeadID, s.ProjectID, s.CreatedAt, s.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("spaceRepo.Create: %w", err)
	}
	if _, err = r.db.Exec(ctx, sql, args...); err != nil {
		if isUniqueViolation(err) {
			return apperr.Conflict("space key already exists")
		}
		return fmt.Errorf("spaceRepo.Create: %w", err)
	}
	return nil
}

func (r *spaceRepo) GetByID(ctx context.Context, id string) (*entity.Space, error) {
	sql, args, err := r.builder.
		Select(spaceCols).From("spaces").
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("spaceRepo.GetByID: %w", err)
	}
	s, err := scanSpace(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("space")
	}
	return s, err
}

func (r *spaceRepo) GetByKey(ctx context.Context, key string) (*entity.Space, error) {
	sql, args, err := r.builder.
		Select(spaceCols).From("spaces").
		Where(sq.And{sq.Eq{"key": key}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("spaceRepo.GetByKey: %w", err)
	}
	s, err := scanSpace(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("space")
	}
	return s, err
}

func (r *spaceRepo) GetByProjectID(ctx context.Context, projectID string) (*entity.Space, error) {
	sql, args, err := r.builder.
		Select(spaceCols).From("spaces").
		Where(sq.And{sq.Eq{"project_id": projectID}, sq.Eq{"deleted_at": nil}}).
		Limit(1).ToSql()
	if err != nil {
		return nil, fmt.Errorf("spaceRepo.GetByProjectID: %w", err)
	}
	s, err := scanSpace(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("space")
	}
	return s, err
}

func (r *spaceRepo) List(ctx context.Context, filter *entity.Filter) ([]*entity.Space, int, error) {
	where := sq.And{sq.Eq{"deleted_at": nil}}
	if filter.Search != "" {
		where = append(where, sq.ILike{"name": "%" + filter.Search + "%"})
	}

	var total int
	cntSQL, cntArgs, _ := r.builder.Select("COUNT(*)").From("spaces").Where(where).ToSql()
	if err := r.db.QueryRow(ctx, cntSQL, cntArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("spaceRepo.List count: %w", err)
	}

	dataSQL, dataArgs, _ := r.builder.
		Select(spaceCols).From("spaces").Where(where).
		OrderBy("created_at DESC").
		Limit(uint64(filter.GetLimit())).Offset(uint64(filter.Offset())).ToSql()

	rows, err := r.db.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("spaceRepo.List query: %w", err)
	}
	defer rows.Close()

	var spaces []*entity.Space
	for rows.Next() {
		s, err := scanSpace(rows)
		if err != nil {
			return nil, 0, err
		}
		spaces = append(spaces, s)
	}
	return spaces, total, rows.Err()
}

func (r *spaceRepo) Update(ctx context.Context, s *entity.Space) error {
	sql, args, err := r.builder.
		Update("spaces").
		Set("name", s.Name).Set("description", s.Description).
		Set("icon_url", s.IconURL).Set("is_archived", s.IsArchived).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": s.ID}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return fmt.Errorf("spaceRepo.Update: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *spaceRepo) SoftDelete(ctx context.Context, id string) error {
	sql, args, err := r.builder.
		Update("spaces").Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return fmt.Errorf("spaceRepo.SoftDelete: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *spaceRepo) ExistsByKey(ctx context.Context, key string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM spaces WHERE key=$1 AND deleted_at IS NULL)`, key,
	).Scan(&exists)
	return exists, err
}

// ─── Members ──────────────────────────────────────────────────────────────────

func (r *spaceRepo) AddMember(ctx context.Context, m *entity.SpaceMember) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO space_members(space_id, user_id, role, created_at) VALUES($1,$2,$3,$4) ON CONFLICT(space_id, user_id) DO UPDATE SET role=EXCLUDED.role`,
		m.SpaceID, m.UserID, m.Role, m.CreatedAt,
	)
	return err
}

func (r *spaceRepo) GetMember(ctx context.Context, spaceID, userID string) (*entity.SpaceMember, error) {
	m := &entity.SpaceMember{User: &entity.UserShort{}}
	err := r.db.QueryRow(ctx, `
		SELECT sm.space_id, sm.user_id, sm.role, sm.created_at,
		       u.id, u.full_name, u.email, u.avatar_url, u.color
		FROM space_members sm
		JOIN users u ON u.id = sm.user_id
		WHERE sm.space_id=$1 AND sm.user_id=$2
	`, spaceID, userID).Scan(
		&m.SpaceID, &m.UserID, &m.Role, &m.CreatedAt,
		&m.User.ID, &m.User.FullName, &m.User.Email, &m.User.AvatarURL, &m.User.Color,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("space member")
	}
	return m, err
}

func (r *spaceRepo) ListMembers(ctx context.Context, spaceID string, filter *entity.Filter) ([]*entity.SpaceMember, int, error) {
	var total int
	if err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM space_members WHERE space_id=$1`, spaceID,
	).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("spaceRepo.ListMembers count: %w", err)
	}

	rows, err := r.db.Query(ctx, `
		SELECT sm.space_id, sm.user_id, sm.role, sm.created_at,
		       u.id, u.full_name, u.email, u.avatar_url, u.color
		FROM space_members sm
		JOIN users u ON u.id = sm.user_id
		WHERE sm.space_id=$1
		ORDER BY sm.created_at ASC
		LIMIT $2 OFFSET $3
	`, spaceID, filter.GetLimit(), filter.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("spaceRepo.ListMembers query: %w", err)
	}
	defer rows.Close()

	var members []*entity.SpaceMember
	for rows.Next() {
		m := &entity.SpaceMember{User: &entity.UserShort{}}
		if err := rows.Scan(
			&m.SpaceID, &m.UserID, &m.Role, &m.CreatedAt,
			&m.User.ID, &m.User.FullName, &m.User.Email, &m.User.AvatarURL, &m.User.Color,
		); err != nil {
			return nil, 0, err
		}
		members = append(members, m)
	}
	return members, total, rows.Err()
}

func (r *spaceRepo) UpdateMemberRole(ctx context.Context, spaceID, userID, role string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE space_members SET role=$3 WHERE space_id=$1 AND user_id=$2`,
		spaceID, userID, role,
	)
	return err
}

func (r *spaceRepo) RemoveMember(ctx context.Context, spaceID, userID string) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM space_members WHERE space_id=$1 AND user_id=$2`, spaceID, userID,
	)
	return err
}

func (r *spaceRepo) IsMember(ctx context.Context, spaceID, userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM space_members WHERE space_id=$1 AND user_id=$2)`,
		spaceID, userID,
	).Scan(&exists)
	return exists, err
}

func (r *spaceRepo) Archive(ctx context.Context, id string) error {
	sql, args, err := r.builder.
		Update("spaces").
		Set("is_archived", true).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).
		ToSql()
	if err != nil {
		return fmt.Errorf("spaceRepo.Archive: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *spaceRepo) Restore(ctx context.Context, id string) error {
	sql, args, err := r.builder.
		Update("spaces").
		Set("is_archived", false).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).
		ToSql()
	if err != nil {
		return fmt.Errorf("spaceRepo.Restore: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *spaceRepo) GetStatistics(ctx context.Context, spaceID string) (*entity.SpaceStatistics, error) {
	stats := &entity.SpaceStatistics{}

	err := r.db.QueryRow(ctx, `
		SELECT
			COUNT(*)                                                        AS total_pages,
			COUNT(*) FILTER (WHERE status = 'published')                   AS published_pages,
			COUNT(*) FILTER (WHERE status = 'draft')                       AS draft_pages
		FROM pages
		WHERE space_id = $1 AND deleted_at IS NULL
	`, spaceID).Scan(&stats.TotalPages, &stats.PublishedPages, &stats.DraftPages)
	if err != nil {
		return nil, fmt.Errorf("spaceRepo.GetStatistics pages: %w", err)
	}

	_ = r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM blog_posts WHERE space_id=$1 AND deleted_at IS NULL`, spaceID,
	).Scan(&stats.TotalBlogPosts)

	_ = r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM space_members WHERE space_id=$1`, spaceID,
	).Scan(&stats.TotalMembers)

	_ = r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM page_views pv
		 JOIN pages p ON p.id = pv.page_id
		 WHERE p.space_id=$1 AND p.deleted_at IS NULL`, spaceID,
	).Scan(&stats.TotalViews)

	_ = r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM pages
		WHERE space_id=$1 AND deleted_at IS NULL
		  AND updated_at >= NOW() - INTERVAL '7 days'
	`, spaceID).Scan(&stats.RecentActivity)

	rows, err := r.db.Query(ctx, `
		SELECT p.author_id, u.full_name, COALESCE(u.avatar_url,''), COUNT(*) AS cnt
		FROM pages p
		JOIN users u ON u.id = p.author_id
		WHERE p.space_id=$1 AND p.deleted_at IS NULL
		GROUP BY p.author_id, u.full_name, u.avatar_url
		ORDER BY cnt DESC
		LIMIT 5
	`, spaceID)
	if err != nil {
		return stats, nil
	}
	defer rows.Close()

	for rows.Next() {
		c := entity.Contributor{}
		if err := rows.Scan(&c.UserID, &c.FullName, &c.AvatarURL, &c.PageCount); err != nil {
			continue
		}
		stats.TopContributors = append(stats.TopContributors, c)
	}
	return stats, rows.Err()
}
