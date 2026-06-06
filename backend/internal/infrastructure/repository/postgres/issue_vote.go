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

type issueVoteRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewIssueVoteRepo(p *pg.Postgres) repository.IssueVoteRepository {
	return &issueVoteRepo{db: p.DB, builder: p.Builder}
}

func (r *issueVoteRepo) Toggle(ctx context.Context, issueID, userID string) (bool, error) {
	var exists bool
	if err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM issue_votes WHERE issue_id=$1 AND user_id=$2)`,
		issueID, userID,
	).Scan(&exists); err != nil {
		return false, fmt.Errorf("issueVoteRepo.Toggle exists: %w", err)
	}

	if exists {
		if _, err := r.db.Exec(ctx,
			`DELETE FROM issue_votes WHERE issue_id=$1 AND user_id=$2`,
			issueID, userID,
		); err != nil {
			return false, fmt.Errorf("issueVoteRepo.Toggle delete: %w", err)
		}
		if _, err := r.db.Exec(ctx,
			`UPDATE issues SET vote_count = (SELECT COUNT(*) FROM issue_votes WHERE issue_id=$1) WHERE id=$1`,
			issueID,
		); err != nil {
			return false, fmt.Errorf("issueVoteRepo.Toggle update count: %w", err)
		}
		return false, nil
	}

	if _, err := r.db.Exec(ctx,
		`INSERT INTO issue_votes(issue_id, user_id, created_at) VALUES($1,$2,$3)`,
		issueID, userID, time.Now().UTC(),
	); err != nil {
		return false, fmt.Errorf("issueVoteRepo.Toggle insert: %w", err)
	}
	if _, err := r.db.Exec(ctx,
		`UPDATE issues SET vote_count = (SELECT COUNT(*) FROM issue_votes WHERE issue_id=$1) WHERE id=$1`,
		issueID,
	); err != nil {
		return false, fmt.Errorf("issueVoteRepo.Toggle update count: %w", err)
	}
	return true, nil
}

func (r *issueVoteRepo) GetSummary(ctx context.Context, issueID, viewerUserID string) (*entity.IssueVoteSummary, error) {
	var summary entity.IssueVoteSummary
	if err := r.db.QueryRow(ctx, `
		SELECT
			COUNT(*),
			COALESCE(BOOL_OR(user_id = $2), FALSE)
		FROM issue_votes
		WHERE issue_id = $1
	`, issueID, viewerUserID).Scan(&summary.Count, &summary.HasMine); err != nil {
		return nil, fmt.Errorf("issueVoteRepo.GetSummary: %w", err)
	}

	rows, err := r.db.Query(ctx, `
		SELECT v.user_id, u.full_name, u.email, u.avatar_url, u.color
		FROM issue_votes v
		JOIN users u ON u.id = v.user_id
		WHERE v.issue_id = $1
		ORDER BY v.created_at ASC
	`, issueID)
	if err != nil {
		return nil, fmt.Errorf("issueVoteRepo.GetSummary voters: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		u := &entity.UserShort{}
		if err := rows.Scan(&u.ID, &u.FullName, &u.Email, &u.AvatarURL, &u.Color); err != nil {
			return nil, err
		}
		summary.Voters = append(summary.Voters, u)
	}
	return &summary, rows.Err()
}
