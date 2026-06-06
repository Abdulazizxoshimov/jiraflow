package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type BoardRepository interface {
	Create(ctx context.Context, b *entity.Board) error
	GetByID(ctx context.Context, id string) (*entity.Board, error)
	GetWithColumns(ctx context.Context, id string) (*entity.Board, error)
	ListByProject(ctx context.Context, projectID string, filter *entity.Filter) ([]*entity.Board, int, error)
	Update(ctx context.Context, b *entity.Board) error
	SoftDelete(ctx context.Context, id string) error

	CreateColumn(ctx context.Context, col *entity.BoardColumn) error
	GetColumnByID(ctx context.Context, id string) (*entity.BoardColumn, error)
	ListColumns(ctx context.Context, boardID string) ([]*entity.BoardColumn, error)
	UpdateColumn(ctx context.Context, col *entity.BoardColumn) error
	DeleteColumn(ctx context.Context, id string) error
	ReorderColumns(ctx context.Context, boardID string, positions map[string]int) error

	SetSwimlaneType(ctx context.Context, boardID, swimlaneType string) error
	GetBoardSwimlanes(ctx context.Context, boardID string, sprintID *string, swimlaneType string) (*entity.GetBoardSwimlanesResp, error)
}
