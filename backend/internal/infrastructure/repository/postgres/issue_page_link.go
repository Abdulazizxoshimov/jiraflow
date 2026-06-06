package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type issuePageLinkRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewIssuePageLinkRepo(p *pg.Postgres) repository.IssuePageLinkRepository {
	return &issuePageLinkRepo{db: p.DB, builder: p.Builder}
}

func (r *issuePageLinkRepo) Create(ctx context.Context, link *entity.IssuePageLink) error {
	sql, args, err := r.builder.
		Insert("issue_page_links").
		Columns("id", "issue_id", "page_id", "linked_by", "created_at").
		Values(link.ID, link.IssueID, link.PageID, link.LinkedBy, link.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("issuePageLinkRepo.Create build: %w", err)
	}
	if _, err = r.db.Exec(ctx, sql, args...); err != nil {
		if isUniqueViolation(err) {
			return apperr.Conflict("issue-page link already exists")
		}
		return fmt.Errorf("issuePageLinkRepo.Create exec: %w", err)
	}
	return nil
}

func (r *issuePageLinkRepo) Delete(ctx context.Context, issueID, pageID string) error {
	sql, args, err := r.builder.
		Delete("issue_page_links").
		Where(sq.Eq{"issue_id": issueID, "page_id": pageID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("issuePageLinkRepo.Delete build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *issuePageLinkRepo) ListByIssue(ctx context.Context, issueID string) ([]*entity.IssuePageLink, error) {
	sql, args, err := r.builder.
		Select("ipl.id", "ipl.issue_id", "ipl.page_id", "ipl.linked_by", "ipl.created_at",
			"p.title", "p.space_id").
		From("issue_page_links ipl").
		Join("pages p ON p.id = ipl.page_id").
		Where(sq.Eq{"ipl.issue_id": issueID}).
		OrderBy("ipl.created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("issuePageLinkRepo.ListByIssue build: %w", err)
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("issuePageLinkRepo.ListByIssue query: %w", err)
	}
	defer rows.Close()

	var links []*entity.IssuePageLink
	for rows.Next() {
		l := &entity.IssuePageLink{Page: &entity.PageShort{}}
		if err := rows.Scan(
			&l.ID, &l.IssueID, &l.PageID, &l.LinkedBy, &l.CreatedAt,
			&l.Page.Title, &l.Page.SpaceID,
		); err != nil {
			return nil, err
		}
		l.Page.ID = l.PageID
		links = append(links, l)
	}
	return links, rows.Err()
}

func (r *issuePageLinkRepo) ListByPage(ctx context.Context, pageID string) ([]*entity.IssuePageLink, error) {
	sql, args, err := r.builder.
		Select("ipl.id", "ipl.issue_id", "ipl.page_id", "ipl.linked_by", "ipl.created_at",
			"i.title", "ws.name", "i.issue_number").
		From("issue_page_links ipl").
		Join("issues i ON i.id = ipl.issue_id").
		Join("workflow_statuses ws ON ws.id = i.status_id").
		Where(sq.Eq{"ipl.page_id": pageID}).
		OrderBy("ipl.created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("issuePageLinkRepo.ListByPage build: %w", err)
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("issuePageLinkRepo.ListByPage query: %w", err)
	}
	defer rows.Close()

	var links []*entity.IssuePageLink
	for rows.Next() {
		var issueNum int
		l := &entity.IssuePageLink{Issue: &entity.IssueShort{}}
		if err := rows.Scan(
			&l.ID, &l.IssueID, &l.PageID, &l.LinkedBy, &l.CreatedAt,
			&l.Issue.Title, &l.Issue.Status, &issueNum,
		); err != nil {
			return nil, err
		}
		l.Issue.ID = l.IssueID
		links = append(links, l)
	}
	return links, rows.Err()
}

func (r *issuePageLinkRepo) Exists(ctx context.Context, issueID, pageID string) (bool, error) {
	sql, args, err := r.builder.
		Select("1").From("issue_page_links").
		Where(sq.Eq{"issue_id": issueID, "page_id": pageID}).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("issuePageLinkRepo.Exists build: %w", err)
	}
	var dummy int
	err = r.db.QueryRow(ctx, sql, args...).Scan(&dummy)
	if err != nil {
		return false, nil
	}
	return true, nil
}
