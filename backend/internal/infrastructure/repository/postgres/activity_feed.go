package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type activityFeedRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewActivityFeedRepo(p *pg.Postgres) repository.ActivityFeedRepository {
	return &activityFeedRepo{db: p.DB, builder: p.Builder}
}

func (r *activityFeedRepo) Create(ctx context.Context, e *entity.ActivityEvent) error {
	metaJSON, err := json.Marshal(e.Meta)
	if err != nil {
		return fmt.Errorf("activityFeedRepo.Create marshal meta: %w", err)
	}
	sql, args, err := r.builder.
		Insert("activity_feed").
		Columns("id", "actor_id", "action", "entity_type", "entity_id", "entity_title",
			"project_id", "space_id", "meta", "created_at").
		Values(e.ID, e.ActorID, e.Action, e.EntityType, e.EntityID, e.EntityTitle,
			e.ProjectID, e.SpaceID, metaJSON, e.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("activityFeedRepo.Create build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *activityFeedRepo) List(ctx context.Context, f *entity.ActivityFilter) ([]*entity.ActivityEvent, int, error) {
	q := r.builder.
		Select("af.id", "af.actor_id", "af.action", "af.entity_type", "af.entity_id",
			"af.entity_title", "af.project_id", "af.space_id", "af.meta", "af.created_at",
			"u.id", "u.full_name", "u.avatar_url", "u.color").
		From("activity_feed af").
		Join("users u ON u.id = af.actor_id").
		OrderBy("af.created_at DESC")

	if f.ActorID != "" {
		q = q.Where(sq.Eq{"af.actor_id": f.ActorID})
	}
	if f.ProjectID != "" {
		q = q.Where(sq.Eq{"af.project_id": f.ProjectID})
	}
	if f.SpaceID != "" {
		q = q.Where(sq.Eq{"af.space_id": f.SpaceID})
	}
	if f.EntityType != "" {
		q = q.Where(sq.Eq{"af.entity_type": f.EntityType})
	}

	countSQL, countArgs, err := r.builder.
		Select("COUNT(*)").From("activity_feed af").ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("activityFeedRepo.List count build: %w", err)
	}
	var total int
	_ = r.db.QueryRow(ctx, countSQL, countArgs...).Scan(&total)

	limit := f.GetLimit()
	offset := f.Offset()
	q = q.Limit(uint64(limit)).Offset(uint64(offset))

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("activityFeedRepo.List build: %w", err)
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("activityFeedRepo.List query: %w", err)
	}
	defer rows.Close()

	var events []*entity.ActivityEvent
	for rows.Next() {
		ev := &entity.ActivityEvent{Actor: &entity.UserShort{}}
		var metaRaw []byte
		if err := rows.Scan(
			&ev.ID, &ev.ActorID, &ev.Action, &ev.EntityType, &ev.EntityID,
			&ev.EntityTitle, &ev.ProjectID, &ev.SpaceID, &metaRaw, &ev.CreatedAt,
			&ev.Actor.ID, &ev.Actor.FullName, &ev.Actor.AvatarURL, &ev.Actor.Color,
		); err != nil {
			return nil, 0, err
		}
		if len(metaRaw) > 0 {
			_ = json.Unmarshal(metaRaw, &ev.Meta)
		}
		events = append(events, ev)
	}
	return events, total, rows.Err()
}
