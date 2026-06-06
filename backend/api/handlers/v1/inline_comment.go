package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

func CreateInlineComment(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageID := c.Param("id")
		authorID := c.GetString(middleware.CtxUserID)
		var req entity.CreateInlineCommentReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		comment, err := h.InlineComment.Create(c.Request.Context(), pageID, authorID, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, comment)
	}
}

func ListInlineComments(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageID := c.Param("id")
		anchorID := c.Query("anchor_id")

		var comments []*entity.InlineComment
		var err error
		if anchorID != "" {
			comments, err = h.InlineComment.ListByAnchor(c.Request.Context(), pageID, anchorID)
		} else {
			comments, err = h.InlineComment.ListByPage(c.Request.Context(), pageID)
		}
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, comments)
	}
}

func UpdateInlineComment(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID := c.GetString(middleware.CtxUserID)
		var req entity.UpdateInlineCommentReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		comment, err := h.InlineComment.Update(c.Request.Context(), c.Param("id"), actorID, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, comment)
	}
}

func ResolveInlineComment(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		resolverID := c.GetString(middleware.CtxUserID)
		if err := h.InlineComment.Resolve(c.Request.Context(), c.Param("id"), resolverID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

func UnresolveInlineComment(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID := c.GetString(middleware.CtxUserID)
		if err := h.InlineComment.Unresolve(c.Request.Context(), c.Param("id"), actorID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

func DeleteInlineComment(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID := c.GetString(middleware.CtxUserID)
		if err := h.InlineComment.Delete(c.Request.Context(), c.Param("id"), actorID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}
