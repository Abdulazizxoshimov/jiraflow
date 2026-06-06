package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type commentRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewCommentRepo(p *pg.Postgres) repository.CommentRepository {
	return &commentRepo{db: p.DB, builder: p.Builder}
}

const commentCols = "id, parent_type, parent_id, author_id, content, content_text, reply_to_id, is_edited, edited_at, created_at, updated_at, deleted_at"

func scanComment(row pgx.Row) (*entity.Comment, error) {
	c := &entity.Comment{}
	var contentJSON []byte
	err := row.Scan(
		&c.ID, &c.ParentType, &c.ParentID, &c.AuthorID,
		&contentJSON, &c.ContentText, &c.ReplyToID,
		&c.IsEdited, &c.EditedAt, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	if len(contentJSON) > 0 {
		_ = json.Unmarshal(contentJSON, &c.Content)
	}
	return c, nil
}

func (r *commentRepo) Create(ctx context.Context, c *entity.Comment) error {
	contentJSON, err := json.Marshal(c.Content)
	if err != nil {
		return fmt.Errorf("commentRepo.Create marshal content: %w", err)
	}
	sql, args, err := r.builder.
		Insert("comments").
		Columns("id", "parent_type", "parent_id", "author_id", "content", "content_text", "reply_to_id", "created_at", "updated_at").
		Values(c.ID, c.ParentType, c.ParentID, c.AuthorID, contentJSON, c.ContentText, c.ReplyToID, c.CreatedAt, c.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("commentRepo.Create: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *commentRepo) GetByID(ctx context.Context, id string) (*entity.Comment, error) {
	sql, args, err := r.builder.
		Select(commentCols).From("comments").
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("commentRepo.GetByID: %w", err)
	}
	c, err := scanComment(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("comment")
	}
	return c, err
}

func (r *commentRepo) ListByParent(ctx context.Context, parentType, parentID string, filter *entity.Filter) ([]*entity.Comment, int, error) {
	where := sq.And{
		sq.Eq{"parent_type": parentType},
		sq.Eq{"parent_id": parentID},
		sq.Eq{"deleted_at": nil},
		sq.Eq{"reply_to_id": nil},
	}

	var total int
	cntSQL, cntArgs, _ := r.builder.Select("COUNT(*)").From("comments").Where(where).ToSql()
	if err := r.db.QueryRow(ctx, cntSQL, cntArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("commentRepo.ListByParent count: %w", err)
	}

	dataSQL, dataArgs, _ := r.builder.
		Select(commentCols).From("comments").Where(where).
		OrderBy("created_at ASC").
		Limit(uint64(filter.GetLimit())).Offset(uint64(filter.Offset())).ToSql()

	rows, err := r.db.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("commentRepo.ListByParent query: %w", err)
	}
	defer rows.Close()

	var comments []*entity.Comment
	for rows.Next() {
		c, err := scanComment(rows)
		if err != nil {
			return nil, 0, err
		}
		comments = append(comments, c)
	}
	return comments, total, rows.Err()
}

func (r *commentRepo) Update(ctx context.Context, c *entity.Comment) error {
	contentJSON, err := json.Marshal(c.Content)
	if err != nil {
		return fmt.Errorf("commentRepo.Update marshal content: %w", err)
	}
	sql, args, err := r.builder.
		Update("comments").
		Set("content", contentJSON).Set("content_text", c.ContentText).
		Set("is_edited", true).Set("edited_at", sq.Expr("NOW()")).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": c.ID}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return fmt.Errorf("commentRepo.Update: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *commentRepo) SoftDelete(ctx context.Context, id string) error {
	sql, args, err := r.builder.
		Update("comments").Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return fmt.Errorf("commentRepo.SoftDelete: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

// ─── Mentions ─────────────────────────────────────────────────────────────────

func (r *commentRepo) AddMention(ctx context.Context, m *entity.CommentMention) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO comment_mentions(comment_id, user_id, created_at) VALUES($1,$2,$3) ON CONFLICT DO NOTHING`,
		m.CommentID, m.UserID, m.CreatedAt,
	)
	return err
}

func (r *commentRepo) ListMentions(ctx context.Context, commentID string) ([]*entity.CommentMention, error) {
	rows, err := r.db.Query(ctx, `
		SELECT cm.comment_id, cm.user_id, cm.created_at,
		       u.id, u.full_name, u.email, u.avatar_url, u.color
		FROM comment_mentions cm
		JOIN users u ON u.id = cm.user_id
		WHERE cm.comment_id = $1
		ORDER BY cm.created_at ASC
	`, commentID)
	if err != nil {
		return nil, fmt.Errorf("commentRepo.ListMentions: %w", err)
	}
	defer rows.Close()

	var mentions []*entity.CommentMention
	for rows.Next() {
		m := &entity.CommentMention{User: &entity.UserShort{}}
		if err := rows.Scan(
			&m.CommentID, &m.UserID, &m.CreatedAt,
			&m.User.ID, &m.User.FullName, &m.User.Email, &m.User.AvatarURL, &m.User.Color,
		); err != nil {
			return nil, err
		}
		mentions = append(mentions, m)
	}
	return mentions, rows.Err()
}

func (r *commentRepo) DeleteMentions(ctx context.Context, commentID string) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM comment_mentions WHERE comment_id=$1`, commentID,
	)
	return err
}

func (r *commentRepo) ToggleReaction(ctx context.Context, commentID, userID, emoji string) error {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM comment_reactions WHERE comment_id=$1 AND user_id=$2 AND emoji=$3)`,
		commentID, userID, emoji,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("commentRepo.ToggleReaction check: %w", err)
	}
	if exists {
		_, err = r.db.Exec(ctx,
			`DELETE FROM comment_reactions WHERE comment_id=$1 AND user_id=$2 AND emoji=$3`,
			commentID, userID, emoji,
		)
	} else {
		_, err = r.db.Exec(ctx,
			`INSERT INTO comment_reactions(id, comment_id, user_id, emoji) VALUES(gen_random_uuid(), $1, $2, $3)`,
			commentID, userID, emoji,
		)
	}
	return err
}

func (r *commentRepo) ListReactions(ctx context.Context, commentID, viewerID string) ([]entity.CommentReactionSummary, error) {
	rows, err := r.db.Query(ctx, `
		SELECT emoji, COUNT(*) AS cnt,
		       BOOL_OR(user_id = $2) AS reacted_by_me
		FROM comment_reactions
		WHERE comment_id = $1
		GROUP BY emoji
		ORDER BY emoji
	`, commentID, viewerID)
	if err != nil {
		return nil, fmt.Errorf("commentRepo.ListReactions: %w", err)
	}
	defer rows.Close()

	var summaries []entity.CommentReactionSummary
	for rows.Next() {
		var s entity.CommentReactionSummary
		if err := rows.Scan(&s.Emoji, &s.Count, &s.ReactedByMe); err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}
	return summaries, rows.Err()
}
