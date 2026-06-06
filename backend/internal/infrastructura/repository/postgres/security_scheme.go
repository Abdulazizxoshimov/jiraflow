package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type securitySchemeRepo struct {
	db *pgxpool.Pool
}

func NewSecuritySchemeRepo(p *pg.Postgres) repository.SecuritySchemeRepository {
	return &securitySchemeRepo{db: p.DB}
}

func (r *securitySchemeRepo) Create(ctx context.Context, req *entity.CreateSecuritySchemeReq) (*entity.SecurityScheme, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	scheme := &entity.SecurityScheme{}
	err = tx.QueryRow(ctx,
		`INSERT INTO security_schemes (name, description, project_id) VALUES ($1,$2,$3)
		 RETURNING id, name, description, project_id, created_at, updated_at`,
		req.Name, req.Description, req.ProjectID,
	).Scan(&scheme.ID, &scheme.Name, &scheme.Description, &scheme.ProjectID, &scheme.CreatedAt, &scheme.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("securitySchemeRepo.Create: %w", err)
	}

	for _, lr := range req.Levels {
		level, err := insertSecurityLevel(ctx, tx, scheme.ID, &lr)
		if err != nil {
			return nil, err
		}
		scheme.Levels = append(scheme.Levels, level)
	}

	return scheme, tx.Commit(ctx)
}

func insertSecurityLevel(ctx context.Context, tx pgx.Tx, schemeID string, req *entity.CreateSecurityLevelReq) (*entity.SecurityLevel, error) {
	level := &entity.SecurityLevel{}
	err := tx.QueryRow(ctx,
		`INSERT INTO security_levels (scheme_id, name, description) VALUES ($1,$2,$3)
		 RETURNING id, scheme_id, name, description, created_at`,
		schemeID, req.Name, req.Description,
	).Scan(&level.ID, &level.SchemeID, &level.Name, &level.Description, &level.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("insertSecurityLevel: %w", err)
	}
	for _, mr := range req.Members {
		m := &entity.SecurityLevelMember{LevelID: level.ID, Type: mr.Type, Value: mr.Value}
		if err := tx.QueryRow(ctx,
			`INSERT INTO security_level_members (level_id, type, value) VALUES ($1,$2,$3) RETURNING id`,
			level.ID, mr.Type, mr.Value,
		).Scan(&m.ID); err != nil {
			return nil, fmt.Errorf("insertSecurityMember: %w", err)
		}
		level.Members = append(level.Members, m)
	}
	return level, nil
}

func (r *securitySchemeRepo) GetByID(ctx context.Context, id string) (*entity.SecurityScheme, error) {
	scheme := &entity.SecurityScheme{}
	err := r.db.QueryRow(ctx,
		`SELECT id, name, description, project_id, created_at, updated_at FROM security_schemes WHERE id=$1`,
		id,
	).Scan(&scheme.ID, &scheme.Name, &scheme.Description, &scheme.ProjectID, &scheme.CreatedAt, &scheme.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("securitySchemeRepo.GetByID: %w", err)
	}
	levels, err := r.loadLevels(ctx, id)
	if err != nil {
		return nil, err
	}
	scheme.Levels = levels
	return scheme, nil
}

func (r *securitySchemeRepo) loadLevels(ctx context.Context, schemeID string) ([]*entity.SecurityLevel, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, scheme_id, name, description, created_at FROM security_levels WHERE scheme_id=$1 ORDER BY created_at`,
		schemeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var levels []*entity.SecurityLevel
	for rows.Next() {
		l := &entity.SecurityLevel{}
		if err := rows.Scan(&l.ID, &l.SchemeID, &l.Name, &l.Description, &l.CreatedAt); err != nil {
			return nil, err
		}
		members, err := r.loadMembers(ctx, l.ID)
		if err != nil {
			return nil, err
		}
		l.Members = members
		levels = append(levels, l)
	}
	return levels, rows.Err()
}

func (r *securitySchemeRepo) loadMembers(ctx context.Context, levelID string) ([]*entity.SecurityLevelMember, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, level_id, type, value FROM security_level_members WHERE level_id=$1`,
		levelID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*entity.SecurityLevelMember
	for rows.Next() {
		m := &entity.SecurityLevelMember{}
		if err := rows.Scan(&m.ID, &m.LevelID, &m.Type, &m.Value); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

func (r *securitySchemeRepo) List(ctx context.Context, projectID string) ([]*entity.SecurityScheme, error) {
	query := `SELECT id, name, description, project_id, created_at, updated_at FROM security_schemes`
	var args []any
	if projectID != "" {
		query += ` WHERE project_id=$1`
		args = append(args, projectID)
	}
	query += ` ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("securitySchemeRepo.List: %w", err)
	}
	defer rows.Close()

	var schemes []*entity.SecurityScheme
	for rows.Next() {
		s := &entity.SecurityScheme{}
		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.ProjectID, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		schemes = append(schemes, s)
	}
	return schemes, rows.Err()
}

func (r *securitySchemeRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM security_schemes WHERE id=$1`, id)
	return err
}

func (r *securitySchemeRepo) AddLevel(ctx context.Context, schemeID string, req *entity.CreateSecurityLevelReq) (*entity.SecurityLevel, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	level, err := insertSecurityLevel(ctx, tx, schemeID, req)
	if err != nil {
		return nil, err
	}
	return level, tx.Commit(ctx)
}

func (r *securitySchemeRepo) GetLevel(ctx context.Context, levelID string) (*entity.SecurityLevel, error) {
	l := &entity.SecurityLevel{}
	err := r.db.QueryRow(ctx,
		`SELECT id, scheme_id, name, description, created_at FROM security_levels WHERE id=$1`,
		levelID,
	).Scan(&l.ID, &l.SchemeID, &l.Name, &l.Description, &l.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("securitySchemeRepo.GetLevel: %w", err)
	}
	members, err := r.loadMembers(ctx, levelID)
	if err != nil {
		return nil, err
	}
	l.Members = members
	return l, nil
}

func (r *securitySchemeRepo) DeleteLevel(ctx context.Context, levelID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM security_levels WHERE id=$1`, levelID)
	return err
}

func (r *securitySchemeRepo) AddMember(ctx context.Context, levelID string, req *entity.CreateSecurityLevelMemberReq) (*entity.SecurityLevelMember, error) {
	m := &entity.SecurityLevelMember{LevelID: levelID, Type: req.Type, Value: req.Value}
	err := r.db.QueryRow(ctx,
		`INSERT INTO security_level_members (level_id, type, value) VALUES ($1,$2,$3) RETURNING id`,
		levelID, req.Type, req.Value,
	).Scan(&m.ID)
	if err != nil {
		return nil, fmt.Errorf("securitySchemeRepo.AddMember: %w", err)
	}
	return m, nil
}

func (r *securitySchemeRepo) DeleteMember(ctx context.Context, memberID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM security_level_members WHERE id=$1`, memberID)
	return err
}
