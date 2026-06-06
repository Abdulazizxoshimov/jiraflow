package postgres

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type pageRestrictionRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewPageRestrictionRepo(p *pg.Postgres) repository.PageRestrictionRepository {
	return &pageRestrictionRepo{db: p.DB, builder: p.Builder}
}

func (r *pageRestrictionRepo) Set(ctx context.Context, pageID string, items []entity.PageRestrictionItem) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pageRestrictionRepo.Set begin: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err = tx.Exec(ctx, `DELETE FROM page_restrictions WHERE page_id=$1`, pageID); err != nil {
		return fmt.Errorf("pageRestrictionRepo.Set delete: %w", err)
	}

	now := time.Now().UTC()
	for _, item := range items {
		if _, err = tx.Exec(ctx,
			`INSERT INTO page_restrictions(id, page_id, type, subject_type, subject_id, created_at)
			 VALUES($1,$2,$3,$4,$5,$6) ON CONFLICT DO NOTHING`,
			uuid.NewString(), pageID, item.Type, item.SubjectType, item.SubjectID, now,
		); err != nil {
			return fmt.Errorf("pageRestrictionRepo.Set insert: %w", err)
		}
	}
	return tx.Commit(ctx)
}

func (r *pageRestrictionRepo) List(ctx context.Context, pageID string) ([]*entity.PageRestriction, error) {
	rows, err := r.db.Query(ctx, `
		SELECT pr.id, pr.page_id, pr.type, pr.subject_type, pr.subject_id, pr.created_at,
		       CASE WHEN pr.subject_type='user' THEN u.id ELSE NULL END,
		       CASE WHEN pr.subject_type='user' THEN u.full_name ELSE NULL END,
		       CASE WHEN pr.subject_type='user' THEN u.email ELSE NULL END,
		       CASE WHEN pr.subject_type='user' THEN u.avatar_url ELSE NULL END,
		       CASE WHEN pr.subject_type='user' THEN u.color ELSE NULL END
		FROM page_restrictions pr
		LEFT JOIN users u ON u.id::text = pr.subject_id AND pr.subject_type='user'
		WHERE pr.page_id = $1
		ORDER BY pr.type, pr.subject_type, pr.created_at
	`, pageID)
	if err != nil {
		return nil, fmt.Errorf("pageRestrictionRepo.List: %w", err)
	}
	defer rows.Close()

	var restrictions []*entity.PageRestriction
	for rows.Next() {
		pr := &entity.PageRestriction{}
		u := &entity.UserShort{}
		if err := rows.Scan(
			&pr.ID, &pr.PageID, &pr.Type, &pr.SubjectType, &pr.SubjectID, &pr.CreatedAt,
			&u.ID, &u.FullName, &u.Email, &u.AvatarURL, &u.Color,
		); err != nil {
			return nil, err
		}
		if u.ID != "" {
			pr.User = u
		}
		restrictions = append(restrictions, pr)
	}
	return restrictions, rows.Err()
}

func (r *pageRestrictionRepo) Clear(ctx context.Context, pageID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM page_restrictions WHERE page_id=$1`, pageID)
	return err
}

func (r *pageRestrictionRepo) CanAccess(ctx context.Context, pageID, userID, accessType string) (bool, error) {
	// Agar restriction yo'q bo'lsa — har kim kira oladi
	var restrictionCount int
	if err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM page_restrictions WHERE page_id=$1 AND type=$2`,
		pageID, accessType,
	).Scan(&restrictionCount); err != nil {
		return false, fmt.Errorf("pageRestrictionRepo.CanAccess count: %w", err)
	}
	if restrictionCount == 0 {
		return true, nil
	}

	// Foydalanuvchi uchun to'g'ridan-to'g'ri ruxsat
	var allowed bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM page_restrictions
			WHERE page_id=$1 AND type=$2 AND subject_type='user' AND subject_id=$3
		)
	`, pageID, accessType, userID).Scan(&allowed)
	return allowed, err
}
