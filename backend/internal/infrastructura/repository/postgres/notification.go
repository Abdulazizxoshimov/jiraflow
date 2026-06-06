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

type notificationRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewNotificationRepo(p *pg.Postgres) repository.NotificationRepository {
	return &notificationRepo{db: p.DB, builder: p.Builder}
}

const notifCols = "id, user_id, type, payload, entity_type, entity_id, actor_id, read_at, email_sent_at, created_at"

func scanNotification(row pgx.Row) (*entity.Notification, error) {
	n := &entity.Notification{}
	var payloadJSON []byte
	err := row.Scan(
		&n.ID, &n.UserID, &n.Type, &payloadJSON,
		&n.EntityType, &n.EntityID, &n.ActorID,
		&n.ReadAt, &n.EmailSentAt, &n.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if len(payloadJSON) > 0 {
		_ = json.Unmarshal(payloadJSON, &n.Payload)
	}
	return n, nil
}

func (r *notificationRepo) Create(ctx context.Context, n *entity.Notification) error {
	payloadJSON, err := json.Marshal(n.Payload)
	if err != nil {
		return fmt.Errorf("notificationRepo.Create marshal payload: %w", err)
	}
	sql, args, err := r.builder.
		Insert("notifications").
		Columns("id", "user_id", "type", "payload", "entity_type", "entity_id", "actor_id", "created_at").
		Values(n.ID, n.UserID, n.Type, payloadJSON, n.EntityType, n.EntityID, n.ActorID, n.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("notificationRepo.Create: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *notificationRepo) GetByID(ctx context.Context, id string) (*entity.Notification, error) {
	sql, args, err := r.builder.
		Select(notifCols).From("notifications").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("notificationRepo.GetByID: %w", err)
	}
	n, err := scanNotification(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("notification")
	}
	return n, err
}

func (r *notificationRepo) ListByUser(ctx context.Context, userID string, filter *entity.NotificationFilter) ([]*entity.Notification, int, error) {
	where := sq.And{sq.Eq{"user_id": userID}}
	if filter.Unread != nil && *filter.Unread {
		where = append(where, sq.Eq{"read_at": nil})
	}

	var total int
	cntSQL, cntArgs, _ := r.builder.Select("COUNT(*)").From("notifications").Where(where).ToSql()
	if err := r.db.QueryRow(ctx, cntSQL, cntArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("notificationRepo.ListByUser count: %w", err)
	}

	dataSQL, dataArgs, _ := r.builder.
		Select(notifCols).From("notifications").Where(where).
		OrderBy("created_at DESC").
		Limit(uint64(filter.GetLimit())).Offset(uint64(filter.Offset())).ToSql()

	rows, err := r.db.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("notificationRepo.ListByUser query: %w", err)
	}
	defer rows.Close()

	var notifs []*entity.Notification
	for rows.Next() {
		n, err := scanNotification(rows)
		if err != nil {
			return nil, 0, err
		}
		notifs = append(notifs, n)
	}
	return notifs, total, rows.Err()
}

func (r *notificationRepo) MarkRead(ctx context.Context, userID string, ids []string) error {
	sql, args, err := r.builder.
		Update("notifications").
		Set("read_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"user_id": userID}, sq.Eq{"id": ids}, sq.Eq{"read_at": nil}}).ToSql()
	if err != nil {
		return fmt.Errorf("notificationRepo.MarkRead: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *notificationRepo) MarkAllRead(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE notifications SET read_at=NOW() WHERE user_id=$1 AND read_at IS NULL`, userID,
	)
	return err
}

func (r *notificationRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM notifications WHERE id=$1`, id)
	return err
}

func (r *notificationRepo) CountUnread(ctx context.Context, userID string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM notifications WHERE user_id=$1 AND read_at IS NULL`, userID,
	).Scan(&count)
	return count, err
}

// ─── Preferences ──────────────────────────────────────────────────────────────

func (r *notificationRepo) GetPreference(ctx context.Context, userID string) (*entity.NotificationPreference, error) {
	pref := &entity.NotificationPreference{}
	err := r.db.QueryRow(ctx, `
		SELECT user_id, email_assigned, email_mentioned, email_commented,
		       email_status, email_watcher, daily_digest, updated_at
		FROM notification_preferences WHERE user_id=$1
	`, userID).Scan(
		&pref.UserID, &pref.EmailAssigned, &pref.EmailMentioned, &pref.EmailCommented,
		&pref.EmailStatus, &pref.EmailWatcher, &pref.DailyDigest, &pref.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return &entity.NotificationPreference{
			UserID:         userID,
			EmailAssigned:  true,
			EmailMentioned: true,
			EmailCommented: true,
			EmailStatus:    true,
			EmailWatcher:   true,
		}, nil
	}
	return pref, err
}

func (r *notificationRepo) UpsertPreference(ctx context.Context, pref *entity.NotificationPreference) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO notification_preferences(user_id, email_assigned, email_mentioned, email_commented, email_status, email_watcher, daily_digest, updated_at)
		VALUES($1,$2,$3,$4,$5,$6,$7,NOW())
		ON CONFLICT(user_id) DO UPDATE SET
			email_assigned=EXCLUDED.email_assigned,
			email_mentioned=EXCLUDED.email_mentioned,
			email_commented=EXCLUDED.email_commented,
			email_status=EXCLUDED.email_status,
			email_watcher=EXCLUDED.email_watcher,
			daily_digest=EXCLUDED.daily_digest,
			updated_at=NOW()
	`, pref.UserID, pref.EmailAssigned, pref.EmailMentioned, pref.EmailCommented,
		pref.EmailStatus, pref.EmailWatcher, pref.DailyDigest,
	)
	return err
}
