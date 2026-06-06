package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// CreateBoard godoc
// @Summary      Create board
// @Tags         boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        project_id  path  string                true  "Project ID"
// @Param        body        body  entity.CreateBoardReq  true  "Board data"
// @Success      201  {object}  object{data=entity.Board}
// @Failure      400  {object}  object{code=string,message=string}
// @Router       /api/v1/projects/{project_id}/boards [post]
func CreateBoard(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		createdBy := c.GetString(middleware.CtxUserID)
		var req entity.CreateBoardReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		board, err := h.Board.Create(c.Request.Context(), projectID, createdBy, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, board)
	}
}

// GetBoard godoc
// @Summary      Get board with columns
// @Tags         boards
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Board ID"
// @Success      200  {object}  object{data=entity.Board}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/boards/{id} [get]
func GetBoard(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		board, err := h.Board.GetWithColumns(c.Request.Context(), id)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, board)
	}
}

// ListBoards godoc
// @Summary      List boards by project
// @Tags         boards
// @Produce      json
// @Security     BearerAuth
// @Param        project_id  path   string  true   "Project ID"
// @Param        page        query  int     false  "Page number"
// @Param        limit       query  int     false  "Page size"
// @Success      200  {object}  object{data=[]entity.Board,total=int,page=int,limit=int,total_pages=int}
// @Router       /api/v1/projects/{project_id}/boards [get]
func ListBoards(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		var filter entity.Filter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		boards, total, err := h.Board.ListByProject(c.Request.Context(), projectID, &filter)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.List(c, boards, total, filter.Page, filter.GetLimit())
	}
}

// UpdateBoard godoc
// @Summary      Update board
// @Tags         boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                true  "Board ID"
// @Param        body  body  entity.UpdateBoardReq  true  "Update data"
// @Success      200  {object}  object{data=entity.Board}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/boards/{id} [put]
func UpdateBoard(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req entity.UpdateBoardReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		board, err := h.Board.Update(c.Request.Context(), id, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, board)
	}
}

// DeleteBoard godoc
// @Summary      Delete board
// @Tags         boards
// @Security     BearerAuth
// @Param        id  path  string  true  "Board ID"
// @Success      204
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/boards/{id} [delete]
func DeleteBoard(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := h.Board.Delete(c.Request.Context(), id); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}
