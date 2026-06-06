package postgres

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type blogPostRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewBlogPostRepo(p *pg.Postgres) repository.BlogPostRepository {
	return &blogPostRepo{db: p.DB, builder: p.Builder}
}

const blogPostCols = `bp.id, bp.space_id, bp.title, bp.body, bp.author_id,
	bp.is_published, bp.published_at, bp.created_at, bp.updated_at,
	u.id, u.full_name, u.email, u.avatar_url, u.color`

func scanBlogPost(row pgx.Row) (*entity.BlogPost, error) {
	bp := &entity.BlogPost{Author: &entity.UserShort{}}
	err := row.Scan(
		&bp.ID, &bp.SpaceID, &bp.Title, &bp.Body, &bp.AuthorID,
		&bp.IsPublished, &bp.PublishedAt, &bp.CreatedAt, &bp.UpdatedAt,
		&bp.Author.ID, &bp.Author.FullName, &bp.Author.Email, &bp.Author.AvatarURL, &bp.Author.Color,
	)
	if err != nil {
		return nil, err
	}
	return bp, nil
}

func (r *blogPostRepo) Create(ctx context.Context, spaceID, authorID string, req *entity.CreateBlogPostReq) (*entity.BlogPost, error) {
	id := uuid.NewString()
	_, err := r.db.Exec(ctx, `
		INSERT INTO blog_posts(id, space_id, title, body, author_id)
		VALUES($1,$2,$3,$4,$5)
	`, id, spaceID, req.Title, req.Body, authorID)
	if err != nil {
		return nil, fmt.Errorf("blogPostRepo.Create: %w", err)
	}
	return r.GetByID(ctx, id)
}

func (r *blogPostRepo) GetByID(ctx context.Context, id string) (*entity.BlogPost, error) {
	row := r.db.QueryRow(ctx, `
		SELECT `+blogPostCols+`
		FROM blog_posts bp
		JOIN users u ON u.id = bp.author_id
		WHERE bp.id = $1 AND bp.deleted_at IS NULL
	`, id)
	bp, err := scanBlogPost(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("blog post")
	}
	return bp, err
}

func (r *blogPostRepo) List(ctx context.Context, filter entity.ListBlogPostsFilter) ([]*entity.BlogPost, int, error) {
	where := sq.And{sq.Eq{"bp.deleted_at": nil}}
	if filter.SpaceID != "" {
		where = append(where, sq.Eq{"bp.space_id": filter.SpaceID})
	}
	if filter.AuthorID != "" {
		where = append(where, sq.Eq{"bp.author_id": filter.AuthorID})
	}
	if filter.OnlyPublished {
		where = append(where, sq.Eq{"bp.is_published": true})
	}

	cntSQL, cntArgs, _ := r.builder.
		Select("COUNT(*)").From("blog_posts bp").Where(where).ToSql()
	var total int
	if err := r.db.QueryRow(ctx, cntSQL, cntArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("blogPostRepo.List count: %w", err)
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}

	dataSQL, dataArgs, _ := r.builder.
		Select(blogPostCols).
		From("blog_posts bp").
		Join("users u ON u.id = bp.author_id").
		Where(where).
		OrderBy("bp.published_at DESC NULLS LAST, bp.created_at DESC").
		Limit(uint64(limit)).Offset(uint64(filter.Offset)).ToSql()

	rows, err := r.db.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("blogPostRepo.List query: %w", err)
	}
	defer rows.Close()

	var posts []*entity.BlogPost
	for rows.Next() {
		bp, err := scanBlogPost(rows)
		if err != nil {
			return nil, 0, err
		}
		posts = append(posts, bp)
	}
	return posts, total, rows.Err()
}

func (r *blogPostRepo) Update(ctx context.Context, id string, req *entity.UpdateBlogPostReq) (*entity.BlogPost, error) {
	q := r.builder.Update("blog_posts").Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}})
	if req.Title != nil {
		q = q.Set("title", *req.Title)
	}
	if req.Body != nil {
		q = q.Set("body", *req.Body)
	}
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("blogPostRepo.Update build: %w", err)
	}
	if _, err := r.db.Exec(ctx, sql, args...); err != nil {
		return nil, fmt.Errorf("blogPostRepo.Update: %w", err)
	}
	return r.GetByID(ctx, id)
}

func (r *blogPostRepo) Delete(ctx context.Context, id string) error {
	sql, args, err := r.builder.
		Update("blog_posts").Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return fmt.Errorf("blogPostRepo.Delete build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *blogPostRepo) Publish(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE blog_posts SET is_published=TRUE, published_at=NOW() WHERE id=$1 AND deleted_at IS NULL`,
		id)
	return err
}

func (r *blogPostRepo) Unpublish(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE blog_posts SET is_published=FALSE, published_at=NULL WHERE id=$1 AND deleted_at IS NULL`,
		id)
	return err
}
