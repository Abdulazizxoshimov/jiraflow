package postgres

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type issuePRRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewIssuePRRepo(p *pg.Postgres) repository.IssuePRRepository {
	return &issuePRRepo{db: p.DB, builder: p.Builder}
}

func (r *issuePRRepo) Create(ctx context.Context, pr *entity.IssuePullRequest) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO issue_pull_requests(id, issue_id, repo_id, pr_number, title, state, url, author_login, created_at, merged_at)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		ON CONFLICT DO NOTHING
	`, pr.ID, pr.IssueID, pr.RepoID, pr.PRNumber, pr.Title, pr.State, pr.URL, pr.AuthorLogin, pr.CreatedAt, pr.MergedAt)
	if err != nil {
		return fmt.Errorf("issuePRRepo.Create: %w", err)
	}
	return nil
}

func (r *issuePRRepo) ListByIssue(ctx context.Context, issueID string) ([]*entity.IssuePullRequest, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, issue_id, repo_id, pr_number, title, state, url, author_login, created_at, merged_at
		FROM issue_pull_requests WHERE issue_id=$1 ORDER BY created_at DESC
	`, issueID)
	if err != nil {
		return nil, fmt.Errorf("issuePRRepo.ListByIssue: %w", err)
	}
	defer rows.Close()
	var prs []*entity.IssuePullRequest
	for rows.Next() {
		pr := &entity.IssuePullRequest{}
		if err := rows.Scan(&pr.ID, &pr.IssueID, &pr.RepoID, &pr.PRNumber, &pr.Title,
			&pr.State, &pr.URL, &pr.AuthorLogin, &pr.CreatedAt, &pr.MergedAt); err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}
	return prs, rows.Err()
}

func (r *issuePRRepo) UpdateState(ctx context.Context, issueID string, prNumber int, state string, mergedAt *string) error {
	var ma *time.Time
	if mergedAt != nil {
		t, err := time.Parse(time.RFC3339, *mergedAt)
		if err == nil {
			ma = &t
		}
	}
	_, err := r.db.Exec(ctx, `
		UPDATE issue_pull_requests SET state=$3, merged_at=$4
		WHERE issue_id=$1 AND pr_number=$2
	`, issueID, prNumber, state, ma)
	return err
}
