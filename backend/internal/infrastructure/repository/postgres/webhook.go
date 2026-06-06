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
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type webhookRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewWebhookRepo(p *pg.Postgres) repository.WebhookRepository {
	return &webhookRepo{db: p.DB, builder: p.Builder}
}

func (r *webhookRepo) Create(ctx context.Context, wh *entity.Webhook) error {
	if wh.ID == "" {
		wh.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	wh.CreatedAt, wh.UpdatedAt = now, now

	_, err := r.db.Exec(ctx,
		`INSERT INTO webhooks(id, project_id, space_id, name, url, secret, events, is_active, created_by, created_at, updated_at)
		 VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		wh.ID, wh.ProjectID, wh.SpaceID, wh.Name, wh.URL, wh.Secret, wh.Events, wh.IsActive, wh.CreatedBy, wh.CreatedAt, wh.UpdatedAt,
	)
	return err
}

func (r *webhookRepo) GetByID(ctx context.Context, id string) (*entity.Webhook, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, project_id, space_id, name, url, secret, events, is_active, created_by, created_at, updated_at
		 FROM webhooks WHERE id=$1`, id,
	)
	wh := &entity.Webhook{}
	err := row.Scan(&wh.ID, &wh.ProjectID, &wh.SpaceID, &wh.Name, &wh.URL, &wh.Secret,
		&wh.Events, &wh.IsActive, &wh.CreatedBy, &wh.CreatedAt, &wh.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("webhook")
	}
	return wh, err
}

func (r *webhookRepo) scanList(ctx context.Context, sql string, args ...any) ([]*entity.Webhook, error) {
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("webhookRepo list: %w", err)
	}
	defer rows.Close()

	var whs []*entity.Webhook
	for rows.Next() {
		wh := &entity.Webhook{}
		if err := rows.Scan(&wh.ID, &wh.ProjectID, &wh.SpaceID, &wh.Name, &wh.URL, &wh.Secret,
			&wh.Events, &wh.IsActive, &wh.CreatedBy, &wh.CreatedAt, &wh.UpdatedAt); err != nil {
			return nil, err
		}
		whs = append(whs, wh)
	}
	return whs, rows.Err()
}

const webhookSelectSQL = `SELECT id, project_id, space_id, name, url, secret, events, is_active, created_by, created_at, updated_at FROM webhooks`

func (r *webhookRepo) ListByProject(ctx context.Context, projectID string) ([]*entity.Webhook, error) {
	return r.scanList(ctx, webhookSelectSQL+` WHERE project_id=$1 ORDER BY created_at DESC`, projectID)
}

func (r *webhookRepo) ListBySpace(ctx context.Context, spaceID string) ([]*entity.Webhook, error) {
	return r.scanList(ctx, webhookSelectSQL+` WHERE space_id=$1 ORDER BY created_at DESC`, spaceID)
}

func (r *webhookRepo) Update(ctx context.Context, wh *entity.Webhook) error {
	_, err := r.db.Exec(ctx,
		`UPDATE webhooks SET name=$1, url=$2, secret=$3, events=$4, is_active=$5, updated_at=NOW()
		 WHERE id=$6`,
		wh.Name, wh.URL, wh.Secret, wh.Events, wh.IsActive, wh.ID,
	)
	return err
}

func (r *webhookRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM webhooks WHERE id=$1`, id)
	return err
}

func (r *webhookRepo) FindByEvent(ctx context.Context, event string, projectID, spaceID *string) ([]*entity.Webhook, error) {
	var args []any
	args = append(args, event)
	where := `$1 = ANY(events) AND is_active = TRUE AND (`
	if projectID != nil {
		where += `project_id=$2`
		args = append(args, *projectID)
	} else if spaceID != nil {
		where += `space_id=$2`
		args = append(args, *spaceID)
	} else {
		return nil, nil
	}
	where += `)`
	return r.scanList(ctx, webhookSelectSQL+` WHERE `+where, args...)
}

func (r *webhookRepo) SaveDelivery(ctx context.Context, d *entity.WebhookDelivery) error {
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	d.DeliveredAt = time.Now().UTC()
	payloadJSON, err := json.Marshal(d.Payload)
	if err != nil {
		return fmt.Errorf("webhookRepo.SaveDelivery marshal: %w", err)
	}
	_, err = r.db.Exec(ctx,
		`INSERT INTO webhook_deliveries(id, webhook_id, event, payload, status_code, response_body, success, delivered_at)
		 VALUES($1,$2,$3,$4,$5,$6,$7,$8)`,
		d.ID, d.WebhookID, d.Event, payloadJSON, d.StatusCode, d.ResponseBody, d.Success, d.DeliveredAt,
	)
	return err
}

func (r *webhookRepo) ListDeliveries(ctx context.Context, webhookID string, limit int) ([]*entity.WebhookDelivery, error) {
	if limit <= 0 {
		limit = 25
	}
	rows, err := r.db.Query(ctx, `
		SELECT id, webhook_id, event, payload, status_code, response_body, success, delivered_at
		FROM webhook_deliveries WHERE webhook_id=$1
		ORDER BY delivered_at DESC LIMIT $2
	`, webhookID, limit)
	if err != nil {
		return nil, fmt.Errorf("webhookRepo.ListDeliveries: %w", err)
	}
	defer rows.Close()

	var deliveries []*entity.WebhookDelivery
	for rows.Next() {
		d := &entity.WebhookDelivery{}
		var payloadJSON []byte
		if err := rows.Scan(&d.ID, &d.WebhookID, &d.Event, &payloadJSON,
			&d.StatusCode, &d.ResponseBody, &d.Success, &d.DeliveredAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(payloadJSON, &d.Payload)
		deliveries = append(deliveries, d)
	}
	return deliveries, rows.Err()
}
