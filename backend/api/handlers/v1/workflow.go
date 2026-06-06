package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// CreateWorkflow godoc
// @Summary      Create workflow
// @Tags         workflows
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  entity.CreateWorkflowReq  true  "Workflow data"
// @Success      201  {object}  object{data=entity.Workflow}
// @Failure      400  {object}  object{code=string,message=string}
// @Router       /api/v1/workflows [post]
func CreateWorkflow(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.CreateWorkflowReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		wf := &entity.Workflow{
			Name:        req.Name,
			Description: req.Description,
			IsDefault:   req.IsDefault,
			CreatedBy:   c.GetString(middleware.CtxUserID),
		}
		created, err := h.Workflow.Create(c.Request.Context(), wf)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, created)
	}
}

// GetWorkflow godoc
// @Summary      Get workflow with statuses and transitions
// @Tags         workflows
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Workflow ID"
// @Success      200  {object}  object{data=entity.Workflow}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/workflows/{id} [get]
func GetWorkflow(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		wf, err := h.Workflow.GetWithDetails(c.Request.Context(), id)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, wf)
	}
}

// ListWorkflows godoc
// @Summary      List workflows
// @Tags         workflows
// @Produce      json
// @Security     BearerAuth
// @Param        page   query  int  false  "Page number"
// @Param        limit  query  int  false  "Page size"
// @Success      200  {object}  object{data=[]entity.Workflow,total=int,page=int,limit=int,total_pages=int}
// @Router       /api/v1/workflows [get]
func ListWorkflows(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter entity.Filter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		wfs, total, err := h.Workflow.List(c.Request.Context(), &filter)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.List(c, wfs, total, filter.Page, filter.GetLimit())
	}
}

// UpdateWorkflow godoc
// @Summary      Update workflow
// @Tags         workflows
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                  true  "Workflow ID"
// @Param        body  body  entity.UpdateWorkflowReq  true  "Update data"
// @Success      200  {object}  object{data=entity.Workflow}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/workflows/{id} [put]
func UpdateWorkflow(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req entity.UpdateWorkflowReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		wf := &entity.Workflow{}
		if req.Name != nil {
			wf.Name = *req.Name
		}
		if req.Description != nil {
			wf.Description = req.Description
		}
		updated, err := h.Workflow.Update(c.Request.Context(), id, wf)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, updated)
	}
}

// DeleteWorkflow godoc
// @Summary      Delete workflow
// @Tags         workflows
// @Security     BearerAuth
// @Param        id  path  string  true  "Workflow ID"
// @Success      204
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/workflows/{id} [delete]
func DeleteWorkflow(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := h.Workflow.Delete(c.Request.Context(), id); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}
