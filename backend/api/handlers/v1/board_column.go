package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// CreateBoardColumn godoc
// @Summary      Add column to board
// @Tags         boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        board_id  path  string                      true  "Board ID"
// @Param        body      body  entity.CreateBoardColumnReq  true  "Column data"
// @Success      201  {object}  object{data=entity.BoardColumn}
// @Failure      400  {object}  object{code=string,message=string}
// @Router       /api/v1/boards/{board_id}/columns [post]
func CreateBoardColumn(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		boardID := c.Param("id")
		var req entity.CreateBoardColumnReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		col, err := h.Board.CreateColumn(c.Request.Context(), boardID, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, col)
	}
}

// UpdateBoardColumn godoc
// @Summary      Update board column
// @Tags         boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                      true  "Column ID"
// @Param        body  body  entity.UpdateBoardColumnReq  true  "Update data"
// @Success      200  {object}  object{data=entity.BoardColumn}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/board-columns/{id} [put]
func UpdateBoardColumn(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req entity.UpdateBoardColumnReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		col, err := h.Board.UpdateColumn(c.Request.Context(), id, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, col)
	}
}

// DeleteBoardColumn godoc
// @Summary      Delete board column
// @Tags         boards
// @Security     BearerAuth
// @Param        id  path  string  true  "Column ID"
// @Success      204
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/board-columns/{id} [delete]
func DeleteBoardColumn(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := h.Board.DeleteColumn(c.Request.Context(), id); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// ReorderBoardColumns godoc
// @Summary      Reorder board columns
// @Tags         boards
// @Accept       json
// @Security     BearerAuth
// @Param        board_id  path  string          true  "Board ID"
// @Param        body      body  object{id=int}  true  "Map of column_id to position"
// @Success      204
// @Router       /api/v1/boards/{board_id}/columns/reorder [put]
func ReorderBoardColumns(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		boardID := c.Param("id")
		var positions map[string]int
		if err := c.ShouldBindJSON(&positions); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		if err := h.Board.ReorderColumns(c.Request.Context(), boardID, positions); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}
