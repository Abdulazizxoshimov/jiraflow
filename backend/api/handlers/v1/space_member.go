package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// AddSpaceMember godoc
// @Summary      Add member to space
// @Tags         space-members
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        space_id  path  string                  true  "Space ID"
// @Param        body      body  entity.AddSpaceMemberReq  true  "Member data"
// @Success      201  {object}  object{data=entity.SpaceMember}
// @Failure      400  {object}  object{code=string,message=string}
// @Router       /api/v1/spaces/{space_id}/members [post]
func AddSpaceMember(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		spaceID := c.Param("id")
		var req entity.AddSpaceMemberReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		m := &entity.SpaceMember{
			SpaceID: spaceID,
			UserID:  req.UserID,
			Role:    req.Role,
		}
		if err := h.Space.AddMember(c.Request.Context(), spaceID, m); err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, m)
	}
}

// ListSpaceMembers godoc
// @Summary      List space members
// @Tags         space-members
// @Produce      json
// @Security     BearerAuth
// @Param        space_id  path   string  true   "Space ID"
// @Param        page      query  int     false  "Page number"
// @Param        limit     query  int     false  "Page size"
// @Success      200  {object}  object{data=[]entity.SpaceMember,total=int,page=int,limit=int,total_pages=int}
// @Router       /api/v1/spaces/{space_id}/members [get]
func ListSpaceMembers(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		spaceID := c.Param("id")
		var filter entity.Filter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		members, total, err := h.Space.ListMembers(c.Request.Context(), spaceID, &filter)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.List(c, members, total, filter.Page, filter.GetLimit())
	}
}

// UpdateSpaceMemberRole godoc
// @Summary      Update space member role
// @Tags         space-members
// @Accept       json
// @Security     BearerAuth
// @Param        space_id  path  string                        true  "Space ID"
// @Param        user_id   path  string                        true  "User ID"
// @Param        body      body  entity.UpdateSpaceMemberRoleReq  true  "Role"
// @Success      204
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/spaces/{space_id}/members/{user_id} [put]
func UpdateSpaceMemberRole(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		spaceID := c.Param("id")
		userID := c.Param("user_id")
		var req entity.UpdateSpaceMemberRoleReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		if err := h.Space.UpdateMemberRole(c.Request.Context(), spaceID, userID, req.Role); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// RemoveSpaceMember godoc
// @Summary      Remove space member
// @Tags         space-members
// @Security     BearerAuth
// @Param        space_id  path  string  true  "Space ID"
// @Param        user_id   path  string  true  "User ID"
// @Success      204
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/spaces/{space_id}/members/{user_id} [delete]
func RemoveSpaceMember(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		spaceID := c.Param("id")
		userID := c.Param("user_id")
		if err := h.Space.RemoveMember(c.Request.Context(), spaceID, userID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}
