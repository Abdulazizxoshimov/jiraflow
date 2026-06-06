package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type inlineCommentRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewInlineCommentRepo(p *pg.Postgres) repository.InlineCommentRepository {
	return &inlineCommentRepo{db: p.DB, builder: p.Builder}
}

const inlineCommentCols = `ic.id, ic.page_id, ic.author_id, ic.anchor_id, ic.quote_text,
	ic.body, ic.resolved, ic.resolved_by, ic.resolved_at, ic.created_at, ic.updated_at`

func scanInlineComment(row pgx.Row) (*entity.InlineComment, error) {
	c := &entity.InlineComment{}
	err := row.Scan(
		&c.ID, &c.PageID, &c.AuthorID, &c.AnchorID, &c.QuoteText,
		&c.Body, &c.Resolved, &c.ResolvedBy, &c.ResolvedAt, &c.CreatedAt, &c.UpdatedAt,
	)
	return c, err
}

func (r *inlineCommentRepo) Create(ctx context.Context, c *entity.InlineComment) error {
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	c.CreatedAt = now
	c.UpdatedAt = now

	_, err := r.db.Exec(ctx,
		`INSERT INTO inline_comments(id, page_id, author_id, anchor_id, quote_text, body, created_at, updated_at)
		 VALUES($1,$2,$3,$4,$5,$6,$7,$8)`,
		c.ID, c.PageID, c.AuthorID, c.AnchorID, c.QuoteText, c.Body, c.CreatedAt, c.UpdatedAt,
	)
	return err
}

func (r *inlineCommentRepo) GetByID(ctx context.Context, id string) (*entity.InlineComment, error) {
	sql, args, err := r.builder.
		Select(inlineCommentCols).From("inline_comments ic").
		Where(sq.And{sq.Eq{"ic.id": id}, sq.Eq{"ic.deleted_at": nil}}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("inlineCommentRepo.GetByID: %w", err)
	}
	c, err := scanInlineComment(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("inline comment")
	}
	return c, err
}

func (r *inlineCommentRepo) listByWhere(ctx context.Context, where sq.Sqlizer) ([]*entity.InlineComment, error) {
	sql, args, err := r.builder.
		Select(inlineCommentCols+`, u.id, u.full_name, u.email, u.avatar_url, u.color`).
		From("inline_comments ic").
		Join("users u ON u.id = ic.author_id").
		Where(where).
		OrderBy("ic.created_at ASC").ToSql()
	if err != nil {
		return nil, fmt.Errorf("inlineCommentRepo list: %w", err)
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("inlineCommentRepo list query: %w", err)
	}
	defer rows.Close()

	var comments []*entity.InlineComment
	for rows.Next() {
		c := &entity.InlineComment{Author: &entity.UserShort{}}
		if err := rows.Scan(
			&c.ID, &c.PageID, &c.AuthorID, &c.AnchorID, &c.QuoteText,
			&c.Body, &c.Resolved, &c.ResolvedBy, &c.ResolvedAt, &c.CreatedAt, &c.UpdatedAt,
			&c.Author.ID, &c.Author.FullName, &c.Author.Email, &c.Author.AvatarURL, &c.Author.Color,
		); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

func (r *inlineCommentRepo) ListByPage(ctx context.Context, pageID string) ([]*entity.InlineComment, error) {
	return r.listByWhere(ctx, sq.And{sq.Eq{"ic.page_id": pageID}, sq.Eq{"ic.deleted_at": nil}})
}

func (r *inlineCommentRepo) ListByAnchor(ctx context.Context, pageID, anchorID string) ([]*entity.InlineComment, error) {
	return r.listByWhere(ctx, sq.And{
		sq.Eq{"ic.page_id": pageID},
		sq.Eq{"ic.anchor_id": anchorID},
		sq.Eq{"ic.deleted_at": nil},
	})
}

func (r *inlineCommentRepo) Update(ctx context.Context, c *entity.InlineComment) error {
	_, err := r.db.Exec(ctx,
		`UPDATE inline_comments SET body=$1, updated_at=NOW() WHERE id=$2 AND deleted_at IS NULL`,
		c.Body, c.ID,
	)
	return err
}

func (r *inlineCommentRepo) Resolve(ctx context.Context, id, resolverID string) error {
	now := time.Now().UTC()
	_, err := r.db.Exec(ctx,
		`UPDATE inline_comments SET resolved=TRUE, resolved_by=$1, resolved_at=$2, updated_at=NOW()
		 WHERE id=$3 AND deleted_at IS NULL`,
		resolverID, now, id,
	)
	return err
}

func (r *inlineCommentRepo) Unresolve(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE inline_comments SET resolved=FALSE, resolved_by=NULL, resolved_at=NULL, updated_at=NOW()
		 WHERE id=$1 AND deleted_at IS NULL`,
		id,
	)
	return err
}

func (r *inlineCommentRepo) SoftDelete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE inline_comments SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id,
	)
	return err
}
