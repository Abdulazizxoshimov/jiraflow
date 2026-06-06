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

type pageReactionRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewPageReactionRepo(p *pg.Postgres) repository.PageReactionRepository {
	return &pageReactionRepo{db: p.DB, builder: p.Builder}
}

// Toggle adds the reaction if it doesn't exist, removes it if it does.
// Returns true if reaction was added, false if removed.
func (r *pageReactionRepo) Toggle(ctx context.Context, pageID, userID, emoji string) (bool, error) {
	var exists bool
	if err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM page_reactions WHERE page_id=$1 AND user_id=$2 AND emoji=$3)`,
		pageID, userID, emoji,
	).Scan(&exists); err != nil {
		return false, fmt.Errorf("pageReactionRepo.Toggle exists: %w", err)
	}

	if exists {
		if _, err := r.db.Exec(ctx,
			`DELETE FROM page_reactions WHERE page_id=$1 AND user_id=$2 AND emoji=$3`,
			pageID, userID, emoji,
		); err != nil {
			return false, fmt.Errorf("pageReactionRepo.Toggle delete: %w", err)
		}
		return false, nil
	}

	if _, err := r.db.Exec(ctx,
		`INSERT INTO page_reactions(id, page_id, user_id, emoji, created_at) VALUES($1,$2,$3,$4,$5)`,
		uuid.NewString(), pageID, userID, emoji, time.Now().UTC(),
	); err != nil {
		return false, fmt.Errorf("pageReactionRepo.Toggle insert: %w", err)
	}
	return true, nil
}

func (r *pageReactionRepo) ListByPage(ctx context.Context, pageID, viewerUserID string) ([]*entity.PageReactionSummary, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			emoji,
			COUNT(*)                                               AS count,
			BOOL_OR(user_id = $2)                                 AS has_mine
		FROM page_reactions
		WHERE page_id = $1
		GROUP BY emoji
		ORDER BY count DESC, emoji ASC
	`, pageID, viewerUserID)
	if err != nil {
		return nil, fmt.Errorf("pageReactionRepo.ListByPage: %w", err)
	}
	defer rows.Close()

	var summaries []*entity.PageReactionSummary
	for rows.Next() {
		s := &entity.PageReactionSummary{}
		if err := rows.Scan(&s.Emoji, &s.Count, &s.HasMine); err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}
	return summaries, rows.Err()
}

func (r *pageReactionRepo) ListUsers(ctx context.Context, pageID, emoji string) ([]*entity.PageReaction, error) {
	rows, err := r.db.Query(ctx, `
		SELECT pr.id, pr.page_id, pr.user_id, pr.emoji, pr.created_at,
		       u.id, u.full_name, u.email, u.avatar_url, u.color
		FROM page_reactions pr
		JOIN users u ON u.id = pr.user_id
		WHERE pr.page_id = $1 AND pr.emoji = $2
		ORDER BY pr.created_at ASC
	`, pageID, emoji)
	if err != nil {
		return nil, fmt.Errorf("pageReactionRepo.ListUsers: %w", err)
	}
	defer rows.Close()

	var reactions []*entity.PageReaction
	for rows.Next() {
		pr := &entity.PageReaction{User: &entity.UserShort{}}
		if err := rows.Scan(
			&pr.ID, &pr.PageID, &pr.UserID, &pr.Emoji, &pr.CreatedAt,
			&pr.User.ID, &pr.User.FullName, &pr.User.Email, &pr.User.AvatarURL, &pr.User.Color,
		); err != nil {
			return nil, err
		}
		reactions = append(reactions, pr)
	}
	return reactions, rows.Err()
}
