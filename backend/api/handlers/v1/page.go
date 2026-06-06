package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// CreatePage godoc
// @Summary      Create page in space
// @Tags         pages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        space_id  path  string                true  "Space ID"
// @Param        body      body  entity.CreatePageReq   true  "Page data"
// @Success      201  {object}  object{data=entity.Page}
// @Failure      400  {object}  object{code=string,message=string}
// @Router       /api/v1/spaces/{space_id}/pages [post]
func CreatePage(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		spaceID := c.Param("id")
		authorID := c.GetString(middleware.CtxUserID)
		var req entity.CreatePageReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		page, err := h.Page.Create(c.Request.Context(), spaceID, authorID, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, page)
	}
}

// GetPage godoc
// @Summary      Get page by ID
// @Tags         pages
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Page ID"
// @Success      200  {object}  object{data=entity.Page}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/pages/{id} [get]
func GetPage(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		page, err := h.Page.GetByID(c.Request.Context(), id)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, page)
	}
}

// ListPages godoc
// @Summary      List pages
// @Tags         pages
// @Produce      json
// @Security     BearerAuth
// @Param        space_id  query  string  false  "Filter by space"
// @Param        parent_id query  string  false  "Filter by parent page"
// @Param        page      query  int     false  "Page number"
// @Param        limit     query  int     false  "Page size"
// @Success      200  {object}  object{data=[]entity.Page,total=int,page=int,limit=int,total_pages=int}
// @Router       /api/v1/pages [get]
func ListPages(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter entity.PageFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		filter.CurrentUserID = c.GetString(middleware.CtxUserID)
		pages, total, err := h.Page.List(c.Request.Context(), &filter)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.List(c, pages, total, filter.Page, filter.GetLimit())
	}
}

// UpdatePage godoc
// @Summary      Update page
// @Tags         pages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                true  "Page ID"
// @Param        body  body  entity.UpdatePageReq   true  "Update data"
// @Success      200  {object}  object{data=entity.Page}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/pages/{id} [put]
func UpdatePage(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		editorID := c.GetString(middleware.CtxUserID)
		var req entity.UpdatePageReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		page, err := h.Page.Update(c.Request.Context(), id, editorID, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, page)
	}
}

// DeletePage godoc
// @Summary      Delete page
// @Tags         pages
// @Security     BearerAuth
// @Param        id  path  string  true  "Page ID"
// @Success      204
// @Failure      403  {object}  object{code=string,message=string}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/pages/{id} [delete]
func DeletePage(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		if err := h.Page.Delete(c.Request.Context(), id, actorID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// GetPageTree godoc
// @Summary      Get page tree for space
// @Tags         pages
// @Produce      json
// @Security     BearerAuth
// @Param        space_id  path  string  true  "Space ID"
// @Success      200  {object}  object{data=[]entity.PageTree}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/spaces/{space_id}/pages/tree [get]
func GetPageTree(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		spaceID := c.Param("id")
		tree, err := h.Page.GetTree(c.Request.Context(), spaceID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, tree)
	}
}

type movePageReq struct {
	Position int     `json:"position"`
	ParentID *string `json:"parent_id"`
}

// MovePage godoc
// @Summary      Move page (change position or parent)
// @Tags         pages
// @Accept       json
// @Security     BearerAuth
// @Param        id    path  string       true  "Page ID"
// @Param        body  body  movePageReq  true  "New position data"
// @Success      204
// @Failure      400  {object}  object{code=string,message=string}
// @Router       /api/v1/pages/{id}/move [put]
func MovePage(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		var req movePageReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		if err := h.Page.Move(c.Request.Context(), id, req.Position, req.ParentID, actorID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// WatchPage godoc
// @Summary      Watch page (subscribe to notifications)
// @Tags         pages
// @Security     BearerAuth
// @Param        id  path  string  true  "Page ID"
// @Success      204
// @Router       /api/v1/pages/{id}/watch [post]
func WatchPage(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID := c.GetString(middleware.CtxUserID)
		if err := h.Page.WatchPage(c.Request.Context(), id, userID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// UnwatchPage godoc
// @Summary      Unwatch page
// @Tags         pages
// @Security     BearerAuth
// @Param        id  path  string  true  "Page ID"
// @Success      204
// @Router       /api/v1/pages/{id}/watch [delete]
func UnwatchPage(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID := c.GetString(middleware.CtxUserID)
		if err := h.Page.UnwatchPage(c.Request.Context(), id, userID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// ListPageWatchers godoc
// @Summary      List page watchers
// @Tags         pages
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Page ID"
// @Success      200  {object}  object{data=[]entity.PageWatcher}
// @Router       /api/v1/pages/{id}/watchers [get]
func ListPageWatchers(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		watchers, err := h.Page.ListWatchers(c.Request.Context(), id)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, watchers)
	}
}

func CopyPage(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.CopyPageReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		newPage, err := h.Page.Copy(c.Request.Context(), c.Param("id"), c.GetString("user_id"), &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, newPage)
	}
}
