package postgres

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type issueAssigneeRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewIssueAssigneeRepo(p *pg.Postgres) repository.IssueAssigneeRepository {
	return &issueAssigneeRepo{db: p.DB, builder: p.Builder}
}

func (r *issueAssigneeRepo) Set(ctx context.Context, issueID string, userIDs []string, primaryID *string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("issueAssigneeRepo.Set begin: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err = tx.Exec(ctx, `DELETE FROM issue_assignees WHERE issue_id=$1`, issueID); err != nil {
		return fmt.Errorf("issueAssigneeRepo.Set delete: %w", err)
	}

	now := time.Now().UTC()
	for _, uid := range userIDs {
		isPrimary := primaryID != nil && *primaryID == uid
		if _, err = tx.Exec(ctx,
			`INSERT INTO issue_assignees(issue_id, user_id, is_primary, created_at)
			 VALUES($1,$2,$3,$4) ON CONFLICT DO NOTHING`,
			issueID, uid, isPrimary, now,
		); err != nil {
			return fmt.Errorf("issueAssigneeRepo.Set insert: %w", err)
		}
	}
	return tx.Commit(ctx)
}

func (r *issueAssigneeRepo) List(ctx context.Context, issueID string) ([]*entity.IssueAssignee, error) {
	rows, err := r.db.Query(ctx, `
		SELECT ia.issue_id, ia.user_id, ia.is_primary, ia.created_at,
		       u.id, u.full_name, u.email, u.avatar_url, u.color
		FROM issue_assignees ia
		JOIN users u ON u.id = ia.user_id
		WHERE ia.issue_id = $1
		ORDER BY ia.is_primary DESC, ia.created_at ASC
	`, issueID)
	if err != nil {
		return nil, fmt.Errorf("issueAssigneeRepo.List: %w", err)
	}
	defer rows.Close()

	var assignees []*entity.IssueAssignee
	for rows.Next() {
		a := &entity.IssueAssignee{User: &entity.UserShort{}}
		if err := rows.Scan(
			&a.IssueID, &a.UserID, &a.IsPrimary, &a.CreatedAt,
			&a.User.ID, &a.User.FullName, &a.User.Email, &a.User.AvatarURL, &a.User.Color,
		); err != nil {
			return nil, err
		}
		assignees = append(assignees, a)
	}
	return assignees, rows.Err()
}

func (r *issueAssigneeRepo) Remove(ctx context.Context, issueID, userID string) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM issue_assignees WHERE issue_id=$1 AND user_id=$2`,
		issueID, userID,
	)
	return err
}
