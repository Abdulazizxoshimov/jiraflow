package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// CreateLabel godoc
// @Summary      Create label
// @Tags         labels
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        project_id  path  string                true  "Project ID"
// @Param        body        body  entity.CreateLabelReq  true  "Label data"
// @Success      201  {object}  object{data=entity.Label}
// @Failure      400  {object}  object{code=string,message=string}
// @Router       /api/v1/projects/{project_id}/labels [post]
func CreateLabel(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		var req entity.CreateLabelReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		label, err := h.Label.Create(c.Request.Context(), projectID, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, label)
	}
}

// GetLabel godoc
// @Summary      Get label by ID
// @Tags         labels
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Label ID"
// @Success      200  {object}  object{data=entity.Label}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/labels/{id} [get]
func GetLabel(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		label, err := h.Label.GetByID(c.Request.Context(), id)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, label)
	}
}

// ListLabels godoc
// @Summary      List labels by project
// @Tags         labels
// @Produce      json
// @Security     BearerAuth
// @Param        project_id  path  string  true  "Project ID"
// @Success      200  {object}  object{data=[]entity.Label}
// @Router       /api/v1/projects/{project_id}/labels [get]
func ListLabels(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		labels, err := h.Label.ListByProject(c.Request.Context(), projectID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, labels)
	}
}

// UpdateLabel godoc
// @Summary      Update label
// @Tags         labels
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                true  "Label ID"
// @Param        body  body  entity.UpdateLabelReq  true  "Update data"
// @Success      200  {object}  object{data=entity.Label}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/labels/{id} [put]
func UpdateLabel(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req entity.UpdateLabelReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		label, err := h.Label.Update(c.Request.Context(), id, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, label)
	}
}

// DeleteLabel godoc
// @Summary      Delete label
// @Tags         labels
// @Security     BearerAuth
// @Param        id  path  string  true  "Label ID"
// @Success      204
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/labels/{id} [delete]
func DeleteLabel(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := h.Label.Delete(c.Request.Context(), id); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}
