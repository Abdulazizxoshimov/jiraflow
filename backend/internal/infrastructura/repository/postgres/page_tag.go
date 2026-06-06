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
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type pageTagRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewPageTagRepo(p *pg.Postgres) repository.PageTagRepository {
	return &pageTagRepo{db: p.DB, builder: p.Builder}
}

func (r *pageTagRepo) Create(ctx context.Context, tag *entity.PageTag) error {
	if tag.ID == "" {
		tag.ID = uuid.NewString()
	}
	if tag.Color == "" {
		tag.Color = "#6B7280"
	}
	tag.CreatedAt = time.Now().UTC()
	sql, args, err := r.builder.
		Insert("page_tags").
		Columns("id", "space_id", "name", "color", "created_at").
		Values(tag.ID, tag.SpaceID, tag.Name, tag.Color, tag.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("pageTagRepo.Create: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *pageTagRepo) GetByID(ctx context.Context, id string) (*entity.PageTag, error) {
	sql, args, err := r.builder.
		Select("id", "space_id", "name", "color", "created_at").
		From("page_tags").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("pageTagRepo.GetByID: %w", err)
	}
	tag := &entity.PageTag{}
	err = r.db.QueryRow(ctx, sql, args...).Scan(&tag.ID, &tag.SpaceID, &tag.Name, &tag.Color, &tag.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("page tag")
	}
	return tag, err
}

func (r *pageTagRepo) List(ctx context.Context, spaceID string) ([]*entity.PageTag, error) {
	sql, args, err := r.builder.
		Select("id", "space_id", "name", "color", "created_at").
		From("page_tags").Where(sq.Eq{"space_id": spaceID}).
		OrderBy("name ASC").ToSql()
	if err != nil {
		return nil, fmt.Errorf("pageTagRepo.List: %w", err)
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("pageTagRepo.List query: %w", err)
	}
	defer rows.Close()

	var tags []*entity.PageTag
	for rows.Next() {
		t := &entity.PageTag{}
		if err := rows.Scan(&t.ID, &t.SpaceID, &t.Name, &t.Color, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}

func (r *pageTagRepo) Update(ctx context.Context, tag *entity.PageTag) error {
	sql, args, err := r.builder.
		Update("page_tags").
		Set("name", tag.Name).Set("color", tag.Color).
		Where(sq.Eq{"id": tag.ID}).ToSql()
	if err != nil {
		return fmt.Errorf("pageTagRepo.Update: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *pageTagRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM page_tags WHERE id=$1`, id)
	return err
}

func (r *pageTagRepo) SetPageTags(ctx context.Context, pageID string, tagIDs []string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pageTagRepo.SetPageTags begin: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err = tx.Exec(ctx, `DELETE FROM page_tag_links WHERE page_id=$1`, pageID); err != nil {
		return fmt.Errorf("pageTagRepo.SetPageTags delete: %w", err)
	}
	for _, tid := range tagIDs {
		if _, err = tx.Exec(ctx,
			`INSERT INTO page_tag_links(page_id, tag_id) VALUES($1,$2) ON CONFLICT DO NOTHING`,
			pageID, tid,
		); err != nil {
			return fmt.Errorf("pageTagRepo.SetPageTags insert: %w", err)
		}
	}
	return tx.Commit(ctx)
}

func (r *pageTagRepo) GetPageTags(ctx context.Context, pageID string) ([]*entity.PageTag, error) {
	rows, err := r.db.Query(ctx, `
		SELECT t.id, t.space_id, t.name, t.color, t.created_at
		FROM page_tags t
		JOIN page_tag_links l ON l.tag_id = t.id
		WHERE l.page_id = $1
		ORDER BY t.name ASC
	`, pageID)
	if err != nil {
		return nil, fmt.Errorf("pageTagRepo.GetPageTags: %w", err)
	}
	defer rows.Close()

	var tags []*entity.PageTag
	for rows.Next() {
		t := &entity.PageTag{}
		if err := rows.Scan(&t.ID, &t.SpaceID, &t.Name, &t.Color, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}

func (r *pageTagRepo) GetPagesByTag(ctx context.Context, tagID string, filter *entity.Filter) ([]*entity.Page, int, error) {
	var total int
	if err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM page_tag_links l
		 JOIN pages p ON p.id = l.page_id
		 WHERE l.tag_id=$1 AND p.deleted_at IS NULL`, tagID,
	).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("pageTagRepo.GetPagesByTag count: %w", err)
	}

	rows, err := r.db.Query(ctx, `
		SELECT p.id, p.space_id, p.parent_id, p.title, p.author_id, p.last_editor_id,
		       p.current_version, p.status, p.position, p.created_at, p.updated_at
		FROM pages p
		JOIN page_tag_links l ON l.page_id = p.id
		WHERE l.tag_id=$1 AND p.deleted_at IS NULL
		ORDER BY p.updated_at DESC
		LIMIT $2 OFFSET $3
	`, tagID, filter.GetLimit(), filter.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("pageTagRepo.GetPagesByTag query: %w", err)
	}
	defer rows.Close()

	var pages []*entity.Page
	for rows.Next() {
		p := &entity.Page{}
		if err := rows.Scan(
			&p.ID, &p.SpaceID, &p.ParentID, &p.Title, &p.AuthorID, &p.LastEditorID,
			&p.CurrentVersion, &p.Status, &p.Position, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		pages = append(pages, p)
	}
	return pages, total, rows.Err()
}
