package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// CreateWorkflowTransition godoc
// @Summary      Add transition to workflow
// @Tags         workflows
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        workflow_id  path  string                              true  "Workflow ID"
// @Param        body         body  entity.CreateWorkflowTransitionReq  true  "Transition data"
// @Success      201  {object}  object{data=entity.WorkflowTransition}
// @Failure      400  {object}  object{code=string,message=string}
// @Router       /api/v1/workflows/{workflow_id}/transitions [post]
func CreateWorkflowTransition(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		workflowID := c.Param("id")
		var req entity.CreateWorkflowTransitionReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		t := &entity.WorkflowTransition{
			WorkflowID:   workflowID,
			FromStatusID: req.FromStatusID,
			ToStatusID:   req.ToStatusID,
			Name:         req.Name,
		}
		created, err := h.Workflow.CreateTransition(c.Request.Context(), t)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, created)
	}
}

// DeleteWorkflowTransition godoc
// @Summary      Delete workflow transition
// @Tags         workflows
// @Security     BearerAuth
// @Param        id  path  string  true  "Transition ID"
// @Success      204
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/workflows/transitions/{id} [delete]
func DeleteWorkflowTransition(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		transitionID := c.Param("transition_id")
		if err := h.Workflow.DeleteTransition(c.Request.Context(), transitionID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}
