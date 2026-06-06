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

type projectRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewProjectRepo(p *pg.Postgres) repository.ProjectRepository {
	return &projectRepo{db: p.DB, builder: p.Builder}
}

const projectCols = "id, key, name, description, icon_url, lead_id, workflow_id, issue_counter, is_archived, created_at, updated_at, deleted_at"

func scanProject(row pgx.Row) (*entity.Project, error) {
	p := &entity.Project{}
	err := row.Scan(
		&p.ID, &p.Key, &p.Name, &p.Description, &p.IconURL,
		&p.LeadID, &p.WorkflowID, &p.IssueCounter, &p.IsArchived,
		&p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
	)
	return p, err
}

func (r *projectRepo) Create(ctx context.Context, p *entity.Project) error {
	sql, args, err := r.builder.
		Insert("projects").
		Columns("id", "key", "name", "description", "icon_url", "lead_id", "workflow_id", "created_at", "updated_at").
		Values(p.ID, p.Key, p.Name, p.Description, p.IconURL, p.LeadID, p.WorkflowID, p.CreatedAt, p.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("projectRepo.Create: %w", err)
	}
	if _, err = r.db.Exec(ctx, sql, args...); err != nil {
		if isUniqueViolation(err) {
			return apperr.Conflict("project key already exists")
		}
		return fmt.Errorf("projectRepo.Create: %w", err)
	}
	return nil
}

func (r *projectRepo) GetByID(ctx context.Context, id string) (*entity.Project, error) {
	sql, args, err := r.builder.
		Select(projectCols).From("projects").
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("projectRepo.GetByID: %w", err)
	}
	p, err := scanProject(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("project")
	}
	return p, err
}

func (r *projectRepo) GetByKey(ctx context.Context, key string) (*entity.Project, error) {
	sql, args, err := r.builder.
		Select(projectCols).From("projects").
		Where(sq.And{sq.Eq{"key": key}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("projectRepo.GetByKey: %w", err)
	}
	p, err := scanProject(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("project")
	}
	return p, err
}

func (r *projectRepo) List(ctx context.Context, filter *entity.ProjectFilter) ([]*entity.Project, int, error) {
	where := sq.And{sq.Eq{"deleted_at": nil}}
	if filter.Search != "" {
		where = append(where, sq.ILike{"name": "%" + filter.Search + "%"})
	}
	if filter.IsArchived != nil {
		where = append(where, sq.Eq{"is_archived": *filter.IsArchived})
	}
	if filter.LeadID != "" {
		where = append(where, sq.Eq{"lead_id": filter.LeadID})
	}

	var total int
	cntSQL, cntArgs, _ := r.builder.Select("COUNT(*)").From("projects").Where(where).ToSql()
	if err := r.db.QueryRow(ctx, cntSQL, cntArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("projectRepo.List count: %w", err)
	}

	order := "DESC"
	if filter.GetSortOrder() == entity.SortAsc {
		order = "ASC"
	}
	dataSQL, dataArgs, _ := r.builder.
		Select(projectCols).From("projects").Where(where).
		OrderBy("created_at " + order).
		Limit(uint64(filter.GetLimit())).Offset(uint64(filter.Offset())).ToSql()

	rows, err := r.db.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("projectRepo.List query: %w", err)
	}
	defer rows.Close()

	var projects []*entity.Project
	for rows.Next() {
		p, err := scanProject(rows)
		if err != nil {
			return nil, 0, err
		}
		projects = append(projects, p)
	}
	return projects, total, rows.Err()
}

func (r *projectRepo) Update(ctx context.Context, p *entity.Project) error {
	sql, args, err := r.builder.
		Update("projects").
		Set("name", p.Name).Set("description", p.Description).
		Set("icon_url", p.IconURL).Set("lead_id", p.LeadID).
		Set("workflow_id", p.WorkflowID).Set("is_archived", p.IsArchived).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": p.ID}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return fmt.Errorf("projectRepo.Update: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *projectRepo) SoftDelete(ctx context.Context, id string) error {
	sql, args, err := r.builder.
		Update("projects").Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return fmt.Errorf("projectRepo.SoftDelete: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *projectRepo) ExistsByKey(ctx context.Context, key string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM projects WHERE key=$1 AND deleted_at IS NULL)`, key,
	).Scan(&exists)
	return exists, err
}

func (r *projectRepo) IncrementIssueCounter(ctx context.Context, id string) (int64, error) {
	var counter int64
	err := r.db.QueryRow(ctx,
		`UPDATE projects SET issue_counter = issue_counter + 1 WHERE id=$1 AND deleted_at IS NULL RETURNING issue_counter`,
		id,
	).Scan(&counter)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, apperr.NotFound("project")
	}
	return counter, err
}

func (r *projectRepo) GetDashboard(ctx context.Context, projectID string) (*entity.ProjectDashboard, error) {
	d := &entity.ProjectDashboard{
		IssuesByPriority: map[string]int{},
		IssuesByType:     map[string]int{},
		IssuesByStatus:   map[string]int{},
	}

	// Open / closed / overdue / total issue counts
	err := r.db.QueryRow(ctx, `
		SELECT
			COUNT(*)                                                              AS total,
			COUNT(*) FILTER (WHERE ws.category != 'done')                        AS open,
			COUNT(*) FILTER (WHERE ws.category = 'done')                         AS closed,
			COUNT(*) FILTER (WHERE ws.category != 'done'
			                   AND i.due_date < NOW())                            AS overdue
		FROM issues i
		JOIN workflow_statuses ws ON ws.id = i.status_id
		WHERE i.project_id = $1 AND i.deleted_at IS NULL
	`, projectID).Scan(&d.TotalIssues, &d.OpenIssues, &d.ClosedIssues, &d.OverdueIssues)
	if err != nil {
		return nil, fmt.Errorf("projectRepo.GetDashboard counts: %w", err)
	}

	// Issues by priority
	rows, err := r.db.Query(ctx,
		`SELECT priority, COUNT(*) FROM issues WHERE project_id=$1 AND deleted_at IS NULL GROUP BY priority`,
		projectID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var priority string
			var cnt int
			if rows.Scan(&priority, &cnt) == nil {
				d.IssuesByPriority[priority] = cnt
			}
		}
	}

	// Issues by type
	rows2, err := r.db.Query(ctx,
		`SELECT type, COUNT(*) FROM issues WHERE project_id=$1 AND deleted_at IS NULL GROUP BY type`,
		projectID)
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var t string
			var cnt int
			if rows2.Scan(&t, &cnt) == nil {
				d.IssuesByType[t] = cnt
			}
		}
	}

	// Issues by status
	rows25, err := r.db.Query(ctx,
		`SELECT ws.name, COUNT(*) FROM issues i JOIN workflow_statuses ws ON ws.id = i.status_id WHERE i.project_id=$1 AND i.deleted_at IS NULL GROUP BY ws.name`,
		projectID)
	if err == nil {
		defer rows25.Close()
		for rows25.Next() {
			var name string
			var cnt int
			if rows25.Scan(&name, &cnt) == nil {
				d.IssuesByStatus[name] = cnt
			}
		}
	}

	// Recent activity (last 7 days)
	_ = r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM issue_history ih JOIN issues i ON i.id = ih.issue_id WHERE i.project_id=$1 AND ih.created_at >= NOW() - INTERVAL '7 days'`,
		projectID).Scan(&d.RecentActivityCount)

	// Top assignees
	rows3, err := r.db.Query(ctx, `
		SELECT i.assignee_id, u.full_name, COUNT(*) AS cnt
		FROM issues i
		JOIN users u ON u.id = i.assignee_id
		WHERE i.project_id=$1 AND i.deleted_at IS NULL AND i.assignee_id IS NOT NULL
		GROUP BY i.assignee_id, u.full_name
		ORDER BY cnt DESC
		LIMIT 10
	`, projectID)
	if err == nil {
		defer rows3.Close()
		for rows3.Next() {
			var a entity.AssigneeDistribution
			if rows3.Scan(&a.UserID, &a.FullName, &a.Count) == nil {
				d.IssuesByAssignee = append(d.IssuesByAssignee, a)
			}
		}
	}

	return d, nil
}
