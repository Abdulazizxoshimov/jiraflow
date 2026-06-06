package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// CreateInvite godoc
// @Summary      Send invitation
// @Tags         invites
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  entity.CreateInviteReq  true  "Invite data"
// @Success      201  {object}  object{data=entity.Invite}
// @Failure      400  {object}  object{code=string,message=string}
// @Router       /api/v1/invites [post]
func CreateInvite(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		invitedBy := c.GetString(middleware.CtxUserID)
		var req entity.CreateInviteReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		invite, err := h.Invite.Create(c.Request.Context(), &req, invitedBy)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, invite)
	}
}

// AcceptInvite godoc
// @Summary      Accept invitation
// @Tags         invites
// @Accept       json
// @Produce      json
// @Param        body  body  entity.AcceptInviteReq  true  "Token and new password"
// @Success      200  {object}  object{data=entity.TokenPair}
// @Failure      400  {object}  object{code=string,message=string}
// @Router       /api/v1/invites/accept [post]
func AcceptInvite(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.AcceptInviteReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		tokens, err := h.Invite.Accept(c.Request.Context(), &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, tokens)
	}
}

// ListPendingInvites godoc
// @Summary      List pending invites
// @Tags         invites
// @Produce      json
// @Security     BearerAuth
// @Param        page   query  int  false  "Page number"
// @Param        limit  query  int  false  "Page size"
// @Success      200  {object}  object{data=[]entity.Invite,total=int,page=int,limit=int,total_pages=int}
// @Router       /api/v1/invites [get]
func ListPendingInvites(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter entity.Filter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		invites, total, err := h.Invite.ListPending(c.Request.Context(), &filter)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.List(c, invites, total, filter.Page, filter.GetLimit())
	}
}

// RevokeInvite godoc
// @Summary      Revoke invitation
// @Tags         invites
// @Security     BearerAuth
// @Param        id  path  string  true  "Invite ID"
// @Success      204
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/invites/{id} [delete]
func RevokeInvite(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := h.Invite.Revoke(c.Request.Context(), id); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}
