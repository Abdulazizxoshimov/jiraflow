package board

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, projectID, createdBy string, req *entity.CreateBoardReq) (*entity.Board, error)
	GetByID(ctx context.Context, id string) (*entity.Board, error)
	GetWithColumns(ctx context.Context, id string) (*entity.Board, error)
	ListByProject(ctx context.Context, projectID string, filter *entity.Filter) ([]*entity.Board, int, error)
	Update(ctx context.Context, id string, req *entity.UpdateBoardReq) (*entity.Board, error)
	Delete(ctx context.Context, id string) error
	CreateColumn(ctx context.Context, boardID string, req *entity.CreateBoardColumnReq) (*entity.BoardColumn, error)
	UpdateColumn(ctx context.Context, id string, req *entity.UpdateBoardColumnReq) (*entity.BoardColumn, error)
	DeleteColumn(ctx context.Context, id string) error
	ReorderColumns(ctx context.Context, boardID string, positions map[string]int) error

	SetSwimlaneType(ctx context.Context, boardID, swimlaneType string) error
	GetBoardSwimlanes(ctx context.Context, boardID string, sprintID *string) (*entity.GetBoardSwimlanesResp, error)
}
