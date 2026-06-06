package postgres

import (
	"context"
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

type componentRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewComponentRepo(p *pg.Postgres) repository.ComponentRepository {
	return &componentRepo{db: p.DB, builder: p.Builder}
}

func scanComponent(row pgx.Row) (*entity.Component, error) {
	c := &entity.Component{}
	var leadID, leadName, leadEmail, leadAvatar *string
	err := row.Scan(
		&c.ID, &c.ProjectID, &c.Name, &c.Description, &c.LeadID,
		&c.CreatedAt, &c.UpdatedAt,
		&leadID, &leadName, &leadEmail, &leadAvatar,
	)
	if err != nil {
		return nil, err
	}
	if leadID != nil {
		c.Lead = &entity.UserShort{ID: *leadID}
		if leadName != nil {
			c.Lead.FullName = *leadName
		}
		if leadEmail != nil {
			c.Lead.Email = *leadEmail
		}
		c.Lead.AvatarURL = leadAvatar
	}
	return c, nil
}

const componentQuery = `
	SELECT c.id, c.project_id, c.name, c.description, c.lead_id, c.created_at, c.updated_at,
	       u.id, u.full_name, u.email, u.avatar_url
	FROM project_components c
	LEFT JOIN users u ON u.id = c.lead_id`

func (r *componentRepo) Create(ctx context.Context, c *entity.Component) error {
	sql, args, err := r.builder.
		Insert("project_components").
		Columns("id", "project_id", "name", "description", "lead_id", "created_at", "updated_at").
		Values(c.ID, c.ProjectID, c.Name, c.Description, c.LeadID, c.CreatedAt, c.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("componentRepo.Create: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *componentRepo) GetByID(ctx context.Context, id string) (*entity.Component, error) {
	query := componentQuery + ` WHERE c.id = $1 AND c.deleted_at IS NULL`
	c, err := scanComponent(r.db.QueryRow(ctx, query, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("component")
	}
	return c, err
}

func (r *componentRepo) List(ctx context.Context, projectID string) ([]*entity.Component, error) {
	query := componentQuery + ` WHERE c.project_id = $1 AND c.deleted_at IS NULL ORDER BY c.name`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("componentRepo.List: %w", err)
	}
	defer rows.Close()

	var result []*entity.Component
	for rows.Next() {
		c, err := scanComponent(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, rows.Err()
}

func (r *componentRepo) Update(ctx context.Context, c *entity.Component) error {
	sql, args, err := r.builder.
		Update("project_components").
		Set("name", c.Name).
		Set("description", c.Description).
		Set("lead_id", c.LeadID).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": c.ID}, sq.Eq{"deleted_at": nil}}).
		ToSql()
	if err != nil {
		return fmt.Errorf("componentRepo.Update: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *componentRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE project_components SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	return err
}

func (r *componentRepo) SetIssueComponents(ctx context.Context, issueID string, componentIDs []string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM issue_components WHERE issue_id=$1`, issueID); err != nil {
		return err
	}
	if len(componentIDs) > 0 {
		b := &pgx.Batch{}
		for _, cid := range componentIDs {
			b.Queue(`INSERT INTO issue_components(issue_id,component_id) VALUES($1,$2) ON CONFLICT DO NOTHING`, issueID, cid)
		}
		br := tx.SendBatch(ctx, b)
		for range componentIDs {
			if _, err := br.Exec(); err != nil {
				br.Close()
				return err
			}
		}
		br.Close()
	}
	return tx.Commit(ctx)
}

func (r *componentRepo) GetIssueComponents(ctx context.Context, issueID string) ([]*entity.Component, error) {
	query := `
		SELECT c.id, c.project_id, c.name, c.description, c.lead_id, c.created_at, c.updated_at,
		       u.id, u.full_name, u.email, u.avatar_url
		FROM project_components c
		JOIN issue_components ic ON ic.component_id = c.id
		LEFT JOIN users u ON u.id = c.lead_id
		WHERE ic.issue_id = $1 AND c.deleted_at IS NULL`
	rows, err := r.db.Query(ctx, query, issueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entity.Component
	for rows.Next() {
		c, err := scanComponent(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, rows.Err()
}
