package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// CreateAutomationRule godoc
// @Summary      Create automation rule for a project
// @Tags         automation
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                          true  "Project ID"
// @Param        body  body  entity.CreateAutomationRuleReq  true  "Rule data"
// @Success      201  {object}  object{data=entity.AutomationRule}
// @Failure      400  {object}  object{code=string,message=string}
// @Router       /api/v1/projects/{id}/automation-rules [post]
func CreateAutomationRule(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		var req entity.CreateAutomationRuleReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		rule, err := h.Automation.Create(c.Request.Context(), projectID, actorID, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, rule)
	}
}

// ListAutomationRules godoc
// @Summary      List automation rules for a project
// @Tags         automation
// @Produce      json
// @Security     BearerAuth
// @Param        id    path   string  true   "Project ID"
// @Param        page  query  int     false  "Page"
// @Param        limit query  int     false  "Limit"
// @Success      200  {object}  object{data=[]entity.AutomationRule,total=int}
// @Router       /api/v1/projects/{id}/automation-rules [get]
func ListAutomationRules(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter entity.AutomationFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		filter.ProjectID = c.Param("id")
		rules, total, err := h.Automation.List(c.Request.Context(), &filter)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.List(c, rules, total, filter.Page, filter.GetLimit())
	}
}

// GetAutomationRule godoc
// @Summary      Get automation rule by ID
// @Tags         automation
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Rule ID"
// @Success      200  {object}  object{data=entity.AutomationRule}
// @Router       /api/v1/automation-rules/{id} [get]
func GetAutomationRule(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		rule, err := h.Automation.GetByID(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, rule)
	}
}

// UpdateAutomationRule godoc
// @Summary      Update automation rule
// @Tags         automation
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                          true  "Rule ID"
// @Param        body  body  entity.UpdateAutomationRuleReq  true  "Update data"
// @Success      200  {object}  object{data=entity.AutomationRule}
// @Router       /api/v1/automation-rules/{id} [put]
func UpdateAutomationRule(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.UpdateAutomationRuleReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		rule, err := h.Automation.Update(c.Request.Context(), c.Param("id"), &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, rule)
	}
}

// DeleteAutomationRule godoc
// @Summary      Delete automation rule
// @Tags         automation
// @Security     BearerAuth
// @Param        id  path  string  true  "Rule ID"
// @Success      204
// @Router       /api/v1/automation-rules/{id} [delete]
func DeleteAutomationRule(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.Automation.Delete(c.Request.Context(), c.Param("id")); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// EnableAutomationRule godoc
// @Summary      Enable automation rule
// @Tags         automation
// @Security     BearerAuth
// @Param        id  path  string  true  "Rule ID"
// @Success      204
// @Router       /api/v1/automation-rules/{id}/enable [post]
func EnableAutomationRule(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.Automation.Enable(c.Request.Context(), c.Param("id")); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// DisableAutomationRule godoc
// @Summary      Disable automation rule
// @Tags         automation
// @Security     BearerAuth
// @Param        id  path  string  true  "Rule ID"
// @Success      204
// @Router       /api/v1/automation-rules/{id}/disable [post]
func DisableAutomationRule(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.Automation.Disable(c.Request.Context(), c.Param("id")); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// ListAutomationLogs godoc
// @Summary      List execution logs for an automation rule
// @Tags         automation
// @Produce      json
// @Security     BearerAuth
// @Param        id    path   string  true   "Rule ID"
// @Param        page  query  int     false  "Page"
// @Param        limit query  int     false  "Limit"
// @Success      200  {object}  object{data=[]entity.AutomationLog,total=int}
// @Router       /api/v1/automation-rules/{id}/logs [get]
func ListAutomationLogs(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter entity.Filter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		logs, total, err := h.Automation.ListLogs(c.Request.Context(), c.Param("id"), &filter)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.List(c, logs, total, filter.Page, filter.GetLimit())
	}
}
