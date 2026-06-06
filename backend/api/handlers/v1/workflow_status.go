package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// CreateWorkflowStatus godoc
// @Summary      Add status to workflow
// @Tags         workflows
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        workflow_id  path  string                        true  "Workflow ID"
// @Param        body         body  entity.CreateWorkflowStatusReq  true  "Status data"
// @Success      201  {object}  object{data=entity.WorkflowStatus}
// @Failure      400  {object}  object{code=string,message=string}
// @Router       /api/v1/workflows/{workflow_id}/statuses [post]
func CreateWorkflowStatus(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		workflowID := c.Param("id")
		var req entity.CreateWorkflowStatusReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		s := &entity.WorkflowStatus{
			WorkflowID: workflowID,
			Name:       req.Name,
			Category:   req.Category,
			Color:      req.Color,
			Position:   req.Position,
			IsInitial:  req.IsInitial,
		}
		status, err := h.Workflow.CreateStatus(c.Request.Context(), s)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, status)
	}
}

// UpdateWorkflowStatus godoc
// @Summary      Update workflow status
// @Tags         workflows
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                        true  "Status ID"
// @Param        body  body  entity.UpdateWorkflowStatusReq  true  "Update data"
// @Success      200  {object}  object{data=entity.WorkflowStatus}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/workflows/statuses/{id} [put]
func UpdateWorkflowStatus(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		statusID := c.Param("id")
		var req entity.UpdateWorkflowStatusReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		s := &entity.WorkflowStatus{ID: statusID}
		if req.Name != nil {
			s.Name = *req.Name
		}
		if req.Category != nil {
			s.Category = *req.Category
		}
		if req.Color != nil {
			s.Color = *req.Color
		}
		status, err := h.Workflow.UpdateStatus(c.Request.Context(), s)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, status)
	}
}

// DeleteWorkflowStatus godoc
// @Summary      Delete workflow status
// @Tags         workflows
// @Security     BearerAuth
// @Param        id  path  string  true  "Status ID"
// @Success      204
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/workflows/statuses/{id} [delete]
func DeleteWorkflowStatus(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		statusID := c.Param("id")
		if err := h.Workflow.DeleteStatus(c.Request.Context(), statusID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}
