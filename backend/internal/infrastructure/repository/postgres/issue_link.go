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

type issueLinkRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewIssueLinkRepo(p *pg.Postgres) repository.IssueLinkRepository {
	return &issueLinkRepo{db: p.DB, builder: p.Builder}
}

func (r *issueLinkRepo) Create(ctx context.Context, link *entity.IssueLink) error {
	sql, args, err := r.builder.
		Insert("issue_links").
		Columns("id", "source_id", "target_id", "link_type", "created_by", "created_at").
		Values(link.ID, link.SourceID, link.TargetID, link.LinkType, link.CreatedBy, link.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("issueLinkRepo.Create: %w", err)
	}
	if _, err = r.db.Exec(ctx, sql, args...); err != nil {
		if isUniqueViolation(err) {
			return apperr.Conflict("issue link already exists")
		}
		return fmt.Errorf("issueLinkRepo.Create: %w", err)
	}
	return nil
}

func (r *issueLinkRepo) GetByID(ctx context.Context, id string) (*entity.IssueLink, error) {
	sql, args, err := r.builder.
		Select("id", "source_id", "target_id", "link_type", "created_by", "created_at").
		From("issue_links").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("issueLinkRepo.GetByID: %w", err)
	}
	link := &entity.IssueLink{}
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&link.ID, &link.SourceID, &link.TargetID, &link.LinkType, &link.CreatedBy, &link.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("issue link")
	}
	return link, err
}

func (r *issueLinkRepo) ListByIssue(ctx context.Context, issueID string) ([]*entity.IssueLink, error) {
	sql, args, err := r.builder.
		Select("id", "source_id", "target_id", "link_type", "created_by", "created_at").
		From("issue_links").
		Where(sq.Or{sq.Eq{"source_id": issueID}, sq.Eq{"target_id": issueID}}).
		OrderBy("created_at DESC").ToSql()
	if err != nil {
		return nil, fmt.Errorf("issueLinkRepo.ListByIssue: %w", err)
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("issueLinkRepo.ListByIssue query: %w", err)
	}
	defer rows.Close()

	var links []*entity.IssueLink
	for rows.Next() {
		link := &entity.IssueLink{}
		if err := rows.Scan(
			&link.ID, &link.SourceID, &link.TargetID, &link.LinkType, &link.CreatedBy, &link.CreatedAt,
		); err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, rows.Err()
}

func (r *issueLinkRepo) Delete(ctx context.Context, id string) error {
	sql, args, err := r.builder.Delete("issue_links").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("issueLinkRepo.Delete: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}
