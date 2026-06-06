package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type issueCommitRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewIssueCommitRepo(p *pg.Postgres) repository.IssueCommitRepository {
	return &issueCommitRepo{db: p.DB, builder: p.Builder}
}

func (r *issueCommitRepo) Create(ctx context.Context, c *entity.IssueCommit) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO issue_commits(id, issue_id, repo_id, sha, message, author_name, author_email, committed_at, url, created_at)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		ON CONFLICT DO NOTHING
	`, c.ID, c.IssueID, c.RepoID, c.SHA, c.Message, c.AuthorName, c.AuthorEmail, c.CommittedAt, c.URL, c.CreatedAt)
	if err != nil {
		return fmt.Errorf("issueCommitRepo.Create: %w", err)
	}
	return nil
}

func (r *issueCommitRepo) ListByIssue(ctx context.Context, issueID string) ([]*entity.IssueCommit, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, issue_id, repo_id, sha, message, author_name, author_email, committed_at, url, created_at
		FROM issue_commits WHERE issue_id=$1 ORDER BY committed_at DESC
	`, issueID)
	if err != nil {
		return nil, fmt.Errorf("issueCommitRepo.ListByIssue: %w", err)
	}
	defer rows.Close()
	var commits []*entity.IssueCommit
	for rows.Next() {
		c := &entity.IssueCommit{}
		if err := rows.Scan(&c.ID, &c.IssueID, &c.RepoID, &c.SHA, &c.Message,
			&c.AuthorName, &c.AuthorEmail, &c.CommittedAt, &c.URL, &c.CreatedAt); err != nil {
			return nil, err
		}
		commits = append(commits, c)
	}
	return commits, rows.Err()
}

func (r *issueCommitRepo) GetBySHA(ctx context.Context, sha string) (*entity.IssueCommit, error) {
	c := &entity.IssueCommit{}
	err := r.db.QueryRow(ctx, `
		SELECT id, issue_id, repo_id, sha, message, author_name, author_email, committed_at, url, created_at
		FROM issue_commits WHERE sha=$1 LIMIT 1
	`, sha).Scan(&c.ID, &c.IssueID, &c.RepoID, &c.SHA, &c.Message,
		&c.AuthorName, &c.AuthorEmail, &c.CommittedAt, &c.URL, &c.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("commit")
	}
	_ = time.Now() // suppress import
	return c, err
}
