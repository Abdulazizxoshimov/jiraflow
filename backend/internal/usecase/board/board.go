package board

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	boardRepo repository.BoardRepository
	log       logger.Logger
}

func New(boardRepo repository.BoardRepository, log logger.Logger) UseCase {
	return &useCase{boardRepo: boardRepo, log: log}
}

func (uc *useCase) Create(ctx context.Context, projectID, createdBy string, req *entity.CreateBoardReq) (*entity.Board, error) {
	now := time.Now().UTC()
	b := &entity.Board{
		ID:        uuid.NewString(),
		ProjectID: projectID,
		Name:      req.Name,
		Type:      req.Type,
		Filter:    req.Filter,
		CreatedBy: createdBy,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if b.Filter == nil {
		b.Filter = map[string]any{}
	}
	if err := uc.boardRepo.Create(ctx, b); err != nil {
		uc.log.Error(ctx, "board.Create: db error", logger.String("project_id", projectID), logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "board created", logger.String("id", b.ID), logger.String("project_id", projectID))
	return b, nil
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.Board, error) {
	return uc.boardRepo.GetByID(ctx, id)
}

func (uc *useCase) GetWithColumns(ctx context.Context, id string) (*entity.Board, error) {
	return uc.boardRepo.GetWithColumns(ctx, id)
}

func (uc *useCase) ListByProject(ctx context.Context, projectID string, filter *entity.Filter) ([]*entity.Board, int, error) {
	return uc.boardRepo.ListByProject(ctx, projectID, filter)
}

func (uc *useCase) Update(ctx context.Context, id string, req *entity.UpdateBoardReq) (*entity.Board, error) {
	b, err := uc.boardRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		b.Name = *req.Name
	}
	if req.Filter != nil {
		b.Filter = req.Filter
	}
	b.UpdatedAt = time.Now().UTC()
	if err := uc.boardRepo.Update(ctx, b); err != nil {
		uc.log.Error(ctx, "board.Update: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "board updated", logger.String("id", id))
	return b, nil
}

func (uc *useCase) Delete(ctx context.Context, id string) error {
	if err := uc.boardRepo.SoftDelete(ctx, id); err != nil {
		uc.log.Error(ctx, "board.Delete: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "board deleted", logger.String("id", id))
	return nil
}

func (uc *useCase) CreateColumn(ctx context.Context, boardID string, req *entity.CreateBoardColumnReq) (*entity.BoardColumn, error) {
	col := &entity.BoardColumn{
		ID:        uuid.NewString(),
		BoardID:   boardID,
		Name:      req.Name,
		Position:  req.Position,
		WIPLimit:  req.WIPLimit,
		StatusIDs: req.StatusIDs,
		CreatedAt: time.Now().UTC(),
	}
	if err := uc.boardRepo.CreateColumn(ctx, col); err != nil {
		uc.log.Error(ctx, "board.CreateColumn: db error", logger.String("board_id", boardID), logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "board column created", logger.String("id", col.ID))
	return col, nil
}

func (uc *useCase) UpdateColumn(ctx context.Context, id string, req *entity.UpdateBoardColumnReq) (*entity.BoardColumn, error) {
	col, err := uc.boardRepo.GetColumnByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		col.Name = *req.Name
	}
	if req.Position != nil {
		col.Position = *req.Position
	}
	if req.WIPLimit != nil {
		col.WIPLimit = req.WIPLimit
	}
	if req.StatusIDs != nil {
		col.StatusIDs = req.StatusIDs
	}
	if err := uc.boardRepo.UpdateColumn(ctx, col); err != nil {
		uc.log.Error(ctx, "board.UpdateColumn: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "board column updated", logger.String("id", id))
	return col, nil
}

func (uc *useCase) DeleteColumn(ctx context.Context, id string) error {
	if _, err := uc.boardRepo.GetColumnByID(ctx, id); err != nil {
		return err
	}
	if err := uc.boardRepo.DeleteColumn(ctx, id); err != nil {
		uc.log.Error(ctx, "board.DeleteColumn: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "board column deleted", logger.String("id", id))
	return nil
}

func (uc *useCase) ReorderColumns(ctx context.Context, boardID string, positions map[string]int) error {
	if _, err := uc.boardRepo.GetByID(ctx, boardID); err != nil {
		return apperr.NotFound("board")
	}
	return uc.boardRepo.ReorderColumns(ctx, boardID, positions)
}

func (uc *useCase) SetSwimlaneType(ctx context.Context, boardID, swimlaneType string) error {
	if _, err := uc.boardRepo.GetByID(ctx, boardID); err != nil {
		return err
	}
	return uc.boardRepo.SetSwimlaneType(ctx, boardID, swimlaneType)
}

func (uc *useCase) GetBoardSwimlanes(ctx context.Context, boardID string, sprintID *string) (*entity.GetBoardSwimlanesResp, error) {
	board, err := uc.boardRepo.GetByID(ctx, boardID)
	if err != nil {
		return nil, err
	}
	return uc.boardRepo.GetBoardSwimlanes(ctx, boardID, sprintID, board.SwimlaneType)
}
