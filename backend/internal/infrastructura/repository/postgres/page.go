package postgres

import (
	"context"
	"encoding/json"
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
	"github.com/jira-backend/jiraflow-backend/internal/pkg/cql"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type pageRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewPageRepo(p *pg.Postgres) repository.PageRepository {
	return &pageRepo{db: p.DB, builder: p.Builder}
}

const pageCols = "id, space_id, parent_id, title, content, content_text, author_id, last_editor_id, current_version, status, position, created_at, updated_at, deleted_at"

func scanPage(row pgx.Row) (*entity.Page, error) {
	p := &entity.Page{}
	var contentJSON []byte
	err := row.Scan(
		&p.ID, &p.SpaceID, &p.ParentID, &p.Title,
		&contentJSON, &p.ContentText, &p.AuthorID, &p.LastEditorID,
		&p.CurrentVersion, &p.Status, &p.Position,
		&p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	if len(contentJSON) > 0 {
		_ = json.Unmarshal(contentJSON, &p.Content)
	}
	return p, nil
}

func (r *pageRepo) Create(ctx context.Context, p *entity.Page) error {
	contentJSON, err := json.Marshal(p.Content)
	if err != nil {
		return fmt.Errorf("pageRepo.Create marshal content: %w", err)
	}
	sql, args, err := r.builder.
		Insert("pages").
		Columns("id", "space_id", "parent_id", "title", "content", "content_text",
			"author_id", "last_editor_id", "current_version", "status", "position", "created_at", "updated_at").
		Values(p.ID, p.SpaceID, p.ParentID, p.Title, contentJSON, p.ContentText,
			p.AuthorID, p.LastEditorID, p.CurrentVersion, p.Status, p.Position, p.CreatedAt, p.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("pageRepo.Create: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *pageRepo) GetByID(ctx context.Context, id string) (*entity.Page, error) {
	sql, args, err := r.builder.
		Select(pageCols).From("pages").
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("pageRepo.GetByID: %w", err)
	}
	p, err := scanPage(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("page")
	}
	return p, err
}

func (r *pageRepo) List(ctx context.Context, filter *entity.PageFilter) ([]*entity.Page, int, error) {
	where := sq.And{sq.Eq{"deleted_at": nil}}
	if filter.SpaceID != "" {
		where = append(where, sq.Eq{"space_id": filter.SpaceID})
	}
	if filter.ParentID != "" {
		where = append(where, sq.Eq{"parent_id": filter.ParentID})
	} else if filter.SpaceID != "" {
		where = append(where, sq.Eq{"parent_id": nil})
	}
	if filter.Status != "" {
		where = append(where, sq.Eq{"status": filter.Status})
	}
	if filter.AuthorID != "" {
		where = append(where, sq.Eq{"author_id": filter.AuthorID})
	}
	if filter.Search != "" {
		where = append(where, sq.ILike{"title": "%" + filter.Search + "%"})
	}
	if filter.CQL != "" {
		if q, err := cql.Parse(filter.CQL); err == nil {
			where = append(where, cql.ToSqlizer(q, filter.CurrentUserID))
		}
	}

	var total int
	cntSQL, cntArgs, _ := r.builder.Select("COUNT(*)").From("pages").Where(where).ToSql()
	if err := r.db.QueryRow(ctx, cntSQL, cntArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("pageRepo.List count: %w", err)
	}

	dataSQL, dataArgs, _ := r.builder.
		Select(pageCols).From("pages").Where(where).
		OrderBy("position ASC, created_at ASC").
		Limit(uint64(filter.GetLimit())).Offset(uint64(filter.Offset())).ToSql()

	rows, err := r.db.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("pageRepo.List query: %w", err)
	}
	defer rows.Close()

	var pages []*entity.Page
	for rows.Next() {
		p, err := scanPage(rows)
		if err != nil {
			return nil, 0, err
		}
		pages = append(pages, p)
	}
	return pages, total, rows.Err()
}

func (r *pageRepo) Update(ctx context.Context, p *entity.Page) error {
	contentJSON, err := json.Marshal(p.Content)
	if err != nil {
		return fmt.Errorf("pageRepo.Update marshal content: %w", err)
	}
	sql, args, err := r.builder.
		Update("pages").
		Set("title", p.Title).Set("content", contentJSON).
		Set("content_text", p.ContentText).Set("status", p.Status).
		Set("last_editor_id", p.LastEditorID).Set("current_version", p.CurrentVersion).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": p.ID}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return fmt.Errorf("pageRepo.Update: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *pageRepo) SoftDelete(ctx context.Context, id string) error {
	sql, args, err := r.builder.
		Update("pages").Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return fmt.Errorf("pageRepo.SoftDelete: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *pageRepo) GetTree(ctx context.Context, spaceID string) ([]*entity.PageTree, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, parent_id, title, position, status
		FROM pages
		WHERE space_id=$1 AND deleted_at IS NULL
		ORDER BY position ASC, created_at ASC
	`, spaceID)
	if err != nil {
		return nil, fmt.Errorf("pageRepo.GetTree: %w", err)
	}
	defer rows.Close()

	all := map[string]*entity.PageTree{}
	var roots []*entity.PageTree

	for rows.Next() {
		n := &entity.PageTree{}
		if err := rows.Scan(&n.ID, &n.ParentID, &n.Title, &n.Position, &n.Status); err != nil {
			return nil, err
		}
		all[n.ID] = n
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, n := range all {
		if n.ParentID == nil {
			roots = append(roots, n)
		} else if parent, ok := all[*n.ParentID]; ok {
			parent.Children = append(parent.Children, *n)
		}
	}
	return roots, nil
}

func (r *pageRepo) UpdatePosition(ctx context.Context, id string, position int, parentID *string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE pages SET position=$2, parent_id=$3, updated_at=NOW() WHERE id=$1 AND deleted_at IS NULL`,
		id, position, parentID,
	)
	return err
}

func (r *pageRepo) GetMaxPosition(ctx context.Context, spaceID string, parentID *string) (int, error) {
	var max int
	err := r.db.QueryRow(ctx,
		`SELECT COALESCE(MAX(position), 0) FROM pages WHERE space_id=$1 AND parent_id IS NOT DISTINCT FROM $2 AND deleted_at IS NULL`,
		spaceID, parentID,
	).Scan(&max)
	return max, err
}

func (r *pageRepo) AddWatcher(ctx context.Context, w *entity.PageWatcher) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO page_watchers (page_id, user_id, created_at) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`,
		w.PageID, w.UserID, w.CreatedAt,
	)
	return err
}

func (r *pageRepo) RemoveWatcher(ctx context.Context, pageID, userID string) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM page_watchers WHERE page_id=$1 AND user_id=$2`,
		pageID, userID,
	)
	return err
}

func (r *pageRepo) ListWatchers(ctx context.Context, pageID string) ([]*entity.PageWatcher, error) {
	rows, err := r.db.Query(ctx,
		`SELECT pw.page_id, pw.user_id, pw.created_at, u.id, u.full_name, u.email, u.avatar_url, u.color
		 FROM page_watchers pw
		 JOIN users u ON u.id = pw.user_id
		 WHERE pw.page_id = $1`,
		pageID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var watchers []*entity.PageWatcher
	for rows.Next() {
		w := &entity.PageWatcher{User: &entity.UserShort{}}
		if err := rows.Scan(&w.PageID, &w.UserID, &w.CreatedAt,
			&w.User.ID, &w.User.FullName, &w.User.Email, &w.User.AvatarURL, &w.User.Color); err != nil {
			return nil, err
		}
		watchers = append(watchers, w)
	}
	return watchers, rows.Err()
}

func (r *pageRepo) IsWatcher(ctx context.Context, pageID, userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM page_watchers WHERE page_id=$1 AND user_id=$2)`,
		pageID, userID,
	).Scan(&exists)
	return exists, err
}

func (r *pageRepo) GetWatcherIDs(ctx context.Context, pageID string) ([]string, error) {
	rows, err := r.db.Query(ctx,
		`SELECT user_id FROM page_watchers WHERE page_id=$1`,
		pageID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// ─── Copy & GetChildren ───────────────────────────────────────────────────────

func (r *pageRepo) Copy(ctx context.Context, srcID, newSpaceID string, newParentID *string, newTitle string, authorID string) (*entity.Page, error) {
	src, err := r.GetByID(ctx, srcID)
	if err != nil {
		return nil, fmt.Errorf("pageRepo.Copy source: %w", err)
	}

	contentJSON, _ := json.Marshal(src.Content)
	now := time.Now().UTC()
	newID := uuid.NewString()

	_, err = r.db.Exec(ctx, `
		INSERT INTO pages(id, space_id, parent_id, title, content, content_text, author_id, last_editor_id, current_version, status, position, created_at, updated_at)
		VALUES($1,$2,$3,$4,$5,$6,$7,$7,1,$8,0,$9,$9)
	`, newID, newSpaceID, newParentID, newTitle, contentJSON, src.ContentText, authorID, src.Status, now)
	if err != nil {
		return nil, fmt.Errorf("pageRepo.Copy insert: %w", err)
	}
	return r.GetByID(ctx, newID)
}

func (r *pageRepo) GetChildren(ctx context.Context, parentID string) ([]*entity.Page, error) {
	sql, args, err := r.builder.
		Select(pageCols).From("pages").
		Where(sq.And{sq.Eq{"parent_id": parentID}, sq.Eq{"deleted_at": nil}}).
		OrderBy("position ASC").ToSql()
	if err != nil {
		return nil, fmt.Errorf("pageRepo.GetChildren build: %w", err)
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("pageRepo.GetChildren query: %w", err)
	}
	defer rows.Close()

	var pages []*entity.Page
	for rows.Next() {
		p, err := scanPage(rows)
		if err != nil {
			return nil, err
		}
		pages = append(pages, p)
	}
	return pages, rows.Err()
}
