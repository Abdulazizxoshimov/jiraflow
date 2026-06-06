package postgres

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type pageViewRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewPageViewRepo(p *pg.Postgres) repository.PageViewRepository {
	return &pageViewRepo{db: p.DB, builder: p.Builder}
}

func (r *pageViewRepo) Record(ctx context.Context, view *entity.PageView) error {
	if view.ID == "" {
		view.ID = uuid.NewString()
	}
	if view.ViewedAt.IsZero() {
		view.ViewedAt = time.Now().UTC()
	}
	_, err := r.db.Exec(ctx,
		`INSERT INTO page_views(id, page_id, user_id, ip_address, viewed_at) VALUES($1,$2,$3,$4,$5)`,
		view.ID, view.PageID, view.UserID, view.IPAddress, view.ViewedAt,
	)
	return err
}

func (r *pageViewRepo) GetAnalytics(ctx context.Context, pageID string) (*entity.PageAnalytics, error) {
	now := time.Now().UTC()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	weekStart := todayStart.AddDate(0, 0, -7)

	a := &entity.PageAnalytics{PageID: pageID}
	err := r.db.QueryRow(ctx, `
		SELECT
			COUNT(*)                                                          AS total_views,
			COUNT(DISTINCT COALESCE(user_id::text, ip_address::text))        AS unique_visitors,
			COUNT(*) FILTER (WHERE viewed_at >= $2)                          AS views_today,
			COUNT(*) FILTER (WHERE viewed_at >= $3)                          AS views_this_week
		FROM page_views
		WHERE page_id = $1
	`, pageID, todayStart, weekStart,
	).Scan(&a.TotalViews, &a.UniqueVisitors, &a.ViewsToday, &a.ViewsThisWeek)
	if err != nil {
		return nil, fmt.Errorf("pageViewRepo.GetAnalytics: %w", err)
	}
	return a, nil
}

func (r *pageViewRepo) ListRecentByUser(ctx context.Context, userID string, limit int) ([]*entity.RecentPage, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.db.Query(ctx, `
		SELECT DISTINCT ON (pv.page_id)
			pv.page_id, p.title, p.space_id, s.name AS space_name, pv.viewed_at
		FROM page_views pv
		JOIN pages p ON p.id = pv.page_id AND p.deleted_at IS NULL
		JOIN spaces s ON s.id = p.space_id AND s.deleted_at IS NULL
		WHERE pv.user_id = $1
		ORDER BY pv.page_id, pv.viewed_at DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("pageViewRepo.ListRecentByUser: %w", err)
	}
	defer rows.Close()

	var recent []*entity.RecentPage
	for rows.Next() {
		rp := &entity.RecentPage{}
		if err := rows.Scan(&rp.PageID, &rp.Title, &rp.SpaceID, &rp.SpaceName, &rp.ViewedAt); err != nil {
			return nil, err
		}
		recent = append(recent, rp)
	}
	return recent, rows.Err()
}
