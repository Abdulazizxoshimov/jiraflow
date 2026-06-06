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

var _ repository.BoardRepository = (*boardRepo)(nil)

type boardRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewBoardRepo(p *pg.Postgres) repository.BoardRepository {
	return &boardRepo{db: p.DB, builder: p.Builder}
}

const boardCols = "id, project_id, name, type, COALESCE(swimlane_type,'none'), filter, created_by, created_at, updated_at, deleted_at"

func scanBoard(row pgx.Row) (*entity.Board, error) {
	b := &entity.Board{}
	var filterJSON []byte
	err := row.Scan(
		&b.ID, &b.ProjectID, &b.Name, &b.Type, &b.SwimlaneType, &filterJSON,
		&b.CreatedBy, &b.CreatedAt, &b.UpdatedAt, &b.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	if len(filterJSON) > 0 {
		_ = json.Unmarshal(filterJSON, &b.Filter)
	}
	return b, nil
}

const boardColCols = "id, board_id, name, position, wip_limit, created_at"

func scanBoardColumn(row pgx.Row) (*entity.BoardColumn, error) {
	col := &entity.BoardColumn{}
	err := row.Scan(&col.ID, &col.BoardID, &col.Name, &col.Position, &col.WIPLimit, &col.CreatedAt)
	return col, err
}

func (r *boardRepo) Create(ctx context.Context, b *entity.Board) error {
	filterJSON, err := json.Marshal(b.Filter)
	if err != nil {
		return fmt.Errorf("boardRepo.Create marshal filter: %w", err)
	}
	sql, args, err := r.builder.
		Insert("boards").
		Columns("id", "project_id", "name", "type", "filter", "created_by", "created_at", "updated_at").
		Values(b.ID, b.ProjectID, b.Name, b.Type, filterJSON, b.CreatedBy, b.CreatedAt, b.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("boardRepo.Create: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *boardRepo) GetByID(ctx context.Context, id string) (*entity.Board, error) {
	sql, args, err := r.builder.
		Select(boardCols).From("boards").
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("boardRepo.GetByID: %w", err)
	}
	b, err := scanBoard(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("board")
	}
	return b, err
}

func (r *boardRepo) GetWithColumns(ctx context.Context, id string) (*entity.Board, error) {
	b, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	cols, err := r.ListColumns(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, col := range cols {
		b.Columns = append(b.Columns, *col)
	}
	return b, nil
}

func (r *boardRepo) ListByProject(ctx context.Context, projectID string, filter *entity.Filter) ([]*entity.Board, int, error) {
	where := sq.And{sq.Eq{"project_id": projectID}, sq.Eq{"deleted_at": nil}}

	var total int
	cntSQL, cntArgs, _ := r.builder.Select("COUNT(*)").From("boards").Where(where).ToSql()
	if err := r.db.QueryRow(ctx, cntSQL, cntArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("boardRepo.ListByProject count: %w", err)
	}

	dataSQL, dataArgs, _ := r.builder.
		Select(boardCols).From("boards").Where(where).
		OrderBy("created_at DESC").
		Limit(uint64(filter.GetLimit())).Offset(uint64(filter.Offset())).ToSql()

	rows, err := r.db.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("boardRepo.ListByProject query: %w", err)
	}
	defer rows.Close()

	var boards []*entity.Board
	for rows.Next() {
		b, err := scanBoard(rows)
		if err != nil {
			return nil, 0, err
		}
		boards = append(boards, b)
	}
	return boards, total, rows.Err()
}

func (r *boardRepo) Update(ctx context.Context, b *entity.Board) error {
	filterJSON, err := json.Marshal(b.Filter)
	if err != nil {
		return fmt.Errorf("boardRepo.Update marshal filter: %w", err)
	}
	sql, args, err := r.builder.
		Update("boards").
		Set("name", b.Name).Set("filter", filterJSON).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": b.ID}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return fmt.Errorf("boardRepo.Update: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *boardRepo) SoftDelete(ctx context.Context, id string) error {
	sql, args, err := r.builder.
		Update("boards").Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return fmt.Errorf("boardRepo.SoftDelete: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

// ─── Columns ──────────────────────────────────────────────────────────────────

func (r *boardRepo) CreateColumn(ctx context.Context, col *entity.BoardColumn) error {
	sql, args, err := r.builder.
		Insert("board_columns").
		Columns("id", "board_id", "name", "position", "wip_limit", "created_at").
		Values(col.ID, col.BoardID, col.Name, col.Position, col.WIPLimit, col.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("boardRepo.CreateColumn: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *boardRepo) GetColumnByID(ctx context.Context, id string) (*entity.BoardColumn, error) {
	sql, args, err := r.builder.
		Select(boardColCols).From("board_columns").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("boardRepo.GetColumnByID: %w", err)
	}
	col, err := scanBoardColumn(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("board column")
	}
	return col, err
}

func (r *boardRepo) ListColumns(ctx context.Context, boardID string) ([]*entity.BoardColumn, error) {
	sql, args, err := r.builder.
		Select(boardColCols).From("board_columns").
		Where(sq.Eq{"board_id": boardID}).
		OrderBy("position ASC").ToSql()
	if err != nil {
		return nil, fmt.Errorf("boardRepo.ListColumns: %w", err)
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("boardRepo.ListColumns query: %w", err)
	}
	defer rows.Close()

	var cols []*entity.BoardColumn
	for rows.Next() {
		col, err := scanBoardColumn(rows)
		if err != nil {
			return nil, err
		}
		cols = append(cols, col)
	}
	return cols, rows.Err()
}

func (r *boardRepo) UpdateColumn(ctx context.Context, col *entity.BoardColumn) error {
	sql, args, err := r.builder.
		Update("board_columns").
		Set("name", col.Name).Set("position", col.Position).Set("wip_limit", col.WIPLimit).
		Where(sq.Eq{"id": col.ID}).ToSql()
	if err != nil {
		return fmt.Errorf("boardRepo.UpdateColumn: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *boardRepo) DeleteColumn(ctx context.Context, id string) error {
	sql, args, err := r.builder.Delete("board_columns").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("boardRepo.DeleteColumn: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *boardRepo) ReorderColumns(ctx context.Context, boardID string, positions map[string]int) error {
	for id, pos := range positions {
		sql, args, err := r.builder.
			Update("board_columns").
			Set("position", pos).
			Where(sq.And{sq.Eq{"id": id}, sq.Eq{"board_id": boardID}}).ToSql()
		if err != nil {
			return fmt.Errorf("boardRepo.ReorderColumns: %w", err)
		}
		if _, err = r.db.Exec(ctx, sql, args...); err != nil {
			return fmt.Errorf("boardRepo.ReorderColumns exec: %w", err)
		}
	}
	return nil
}

// ─── Swimlanes ─────────────────────────────────────────────────────────────────

func (r *boardRepo) SetSwimlaneType(ctx context.Context, boardID, swimlaneType string) error {
	sql, args, err := r.builder.
		Update("boards").
		Set("swimlane_type", swimlaneType).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": boardID}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return fmt.Errorf("boardRepo.SetSwimlaneType: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *boardRepo) GetBoardSwimlanes(ctx context.Context, boardID string, sprintID *string, swimlaneType string) (*entity.GetBoardSwimlanesResp, error) {
	board, err := r.GetByID(ctx, boardID)
	if err != nil {
		return nil, err
	}
	cols, err := r.ListColumns(ctx, boardID)
	if err != nil {
		return nil, err
	}

	// fetch issues for this board's project (optionally sprint-filtered)
	issueQuery := r.builder.
		Select(`i.id, i.project_id, i.issue_number, i.title, i.description,
			i.type, i.status_id, i.priority, i.assignee_id, i.reporter_id,
			i.parent_id, i.sprint_id, i.story_points, i.due_date,
			i.original_estimate, i.remaining_estimate, i.custom_fields,
			i.position, i.created_at, i.updated_at, i.deleted_at,
			u.id, u.full_name, u.email, u.avatar_url`).
		From("issues i").
		LeftJoin("users u ON u.id = i.assignee_id").
		Where(sq.And{sq.Eq{"i.project_id": board.ProjectID}, sq.Eq{"i.deleted_at": nil}})

	if sprintID != nil {
		issueQuery = issueQuery.Where(sq.Eq{"i.sprint_id": *sprintID})
	}

	issueSql, issueArgs, err := issueQuery.OrderBy("i.position ASC, i.created_at ASC").ToSql()
	if err != nil {
		return nil, fmt.Errorf("boardRepo.GetBoardSwimlanes query build: %w", err)
	}

	rows, err := r.db.Query(ctx, issueSql, issueArgs...)
	if err != nil {
		return nil, fmt.Errorf("boardRepo.GetBoardSwimlanes query: %w", err)
	}
	defer rows.Close()

	type issueWithAssignee struct {
		issue    *entity.Issue
		assignee *entity.UserShort
	}
	var issues []issueWithAssignee
	for rows.Next() {
		i := &entity.Issue{}
		var cfJSON []byte
		a := &entity.UserShort{}
		var aID, aName, aEmail, aAvatar *string
		if err := rows.Scan(
			&i.ID, &i.ProjectID, &i.IssueNumber, &i.Title, &i.Description,
			&i.Type, &i.StatusID, &i.Priority, &i.AssigneeID, &i.ReporterID,
			&i.ParentID, &i.SprintID, &i.StoryPoints, &i.DueDate,
			&i.OriginalEstimate, &i.RemainingEstimate, &cfJSON,
			&i.Position, &i.CreatedAt, &i.UpdatedAt, &i.DeletedAt,
			&aID, &aName, &aEmail, &aAvatar,
		); err != nil {
			return nil, fmt.Errorf("boardRepo.GetBoardSwimlanes scan: %w", err)
		}
		if len(cfJSON) > 0 {
			_ = json.Unmarshal(cfJSON, &i.CustomFields)
		}
		if aID != nil {
			a.ID = *aID
			if aName != nil {
				a.FullName = *aName
			}
			if aEmail != nil {
				a.Email = *aEmail
			}
			a.AvatarURL = aAvatar
			i.Assignee = a
		}
		issues = append(issues, issueWithAssignee{issue: i, assignee: a})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// group issues into swimlanes in Go
	type key struct{ k, label string }
	order := []key{}
	groups := map[string]*entity.BoardSwimlane{}

	add := func(k, label string, issue *entity.Issue) {
		if _, ok := groups[k]; !ok {
			groups[k] = &entity.BoardSwimlane{Key: k, Label: label}
			order = append(order, key{k, label})
		}
		groups[k].Issues = append(groups[k].Issues, issue)
	}

	for _, iwa := range issues {
		i := iwa.issue
		switch swimlaneType {
		case "assignee":
			if i.AssigneeID != nil && i.Assignee != nil {
				add(*i.AssigneeID, i.Assignee.FullName, i)
			} else {
				add("unassigned", "Unassigned", i)
			}
		case "epic":
			if i.ParentID != nil {
				add(*i.ParentID, "Epic "+(*i.ParentID)[:8], i)
			} else {
				add("no_epic", "No Epic", i)
			}
		case "priority":
			add(i.Priority, i.Priority, i)
		case "label":
			// without label join: put in default
			add("default", "All Issues", i)
		default:
			add("default", "All Issues", i)
		}
	}

	swimlanes := make([]*entity.BoardSwimlane, 0, len(order))
	for _, o := range order {
		swimlanes = append(swimlanes, groups[o.k])
	}

	return &entity.GetBoardSwimlanesResp{
		SwimlaneType: board.SwimlaneType,
		Columns:      cols,
		Swimlanes:    swimlanes,
	}, nil
}
