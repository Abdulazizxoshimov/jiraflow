package v1

import (
	"strings"

	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// normalizeParentType "issues" → "issue", "pages" → "page"
func normalizeParentType(t string) string {
	return strings.TrimSuffix(t, "s")
}

// CreateComment godoc
// @Summary      Create comment
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        parent_type  path  string                  true  "Parent type (issues|pages)"
// @Param        parent_id    path  string                  true  "Parent ID"
// @Param        body         body  entity.CreateCommentReq  true  "Comment data"
// @Success      201  {object}  object{data=entity.Comment}
// @Failure      400  {object}  object{code=string,message=string}
// @Router       /api/v1/{parent_type}/{parent_id}/comments [post]
func CreateComment(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		parentType := normalizeParentType(c.Param("parent_type"))
		parentID := c.Param("parent_id")
		authorID := c.GetString(middleware.CtxUserID)
		var req entity.CreateCommentReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		comment, err := h.Comment.Create(c.Request.Context(), parentType, parentID, authorID, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, comment)
	}
}

// GetComment godoc
// @Summary      Get comment by ID
// @Tags         comments
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Comment ID"
// @Success      200  {object}  object{data=entity.Comment}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/comments/{id} [get]
func GetComment(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		comment, err := h.Comment.GetByID(c.Request.Context(), id)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, comment)
	}
}

// ListComments godoc
// @Summary      List comments for a parent
// @Tags         comments
// @Produce      json
// @Security     BearerAuth
// @Param        parent_type  path   string  true   "Parent type (issues|pages)"
// @Param        parent_id    path   string  true   "Parent ID"
// @Param        page         query  int     false  "Page number"
// @Param        limit        query  int     false  "Page size"
// @Success      200  {object}  object{data=[]entity.Comment,total=int,page=int,limit=int,total_pages=int}
// @Router       /api/v1/{parent_type}/{parent_id}/comments [get]
func ListComments(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		parentType := normalizeParentType(c.Param("parent_type"))
		parentID := c.Param("parent_id")
		var filter entity.Filter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		comments, total, err := h.Comment.ListByParent(c.Request.Context(), parentType, parentID, &filter)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.List(c, comments, total, filter.Page, filter.GetLimit())
	}
}

// UpdateComment godoc
// @Summary      Update comment
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                  true  "Comment ID"
// @Param        body  body  entity.UpdateCommentReq  true  "Update data"
// @Success      200  {object}  object{data=entity.Comment}
// @Failure      403  {object}  object{code=string,message=string}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/comments/{id} [put]
func UpdateComment(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		var req entity.UpdateCommentReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		comment, err := h.Comment.Update(c.Request.Context(), id, actorID, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, comment)
	}
}

// DeleteComment godoc
// @Summary      Delete comment
// @Tags         comments
// @Security     BearerAuth
// @Param        id  path  string  true  "Comment ID"
// @Success      204
// @Failure      403  {object}  object{code=string,message=string}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/comments/{id} [delete]
func DeleteComment(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		if err := h.Comment.Delete(c.Request.Context(), id, actorID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// ToggleCommentReaction godoc
// @Summary      Toggle emoji reaction on a comment
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string  true  "Comment ID"
// @Param        body  body  object{emoji=string}  true  "Emoji"
// @Success      204
// @Router       /api/v1/comments/{id}/reactions [post]
func ToggleCommentReaction(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID := c.GetString(middleware.CtxUserID)
		var req struct {
			Emoji string `json:"emoji" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		if err := h.Comment.ToggleReaction(c.Request.Context(), id, userID, req.Emoji); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// ListCommentReactions godoc
// @Summary      List reactions for a comment
// @Tags         comments
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Comment ID"
// @Success      200  {object}  object{data=[]entity.CommentReactionSummary}
// @Router       /api/v1/comments/{id}/reactions [get]
func ListCommentReactions(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		viewerID := c.GetString(middleware.CtxUserID)
		reactions, err := h.Comment.ListReactions(c.Request.Context(), id, viewerID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, reactions)
	}
}
