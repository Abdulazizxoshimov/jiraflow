package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// CreateAPIKey godoc
// @Summary      Create a new API key
// @Tags         api-keys
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  entity.CreateAPIKeyReq  true  "API key params"
// @Success      201  {object}  object{data=entity.CreateAPIKeyResp}
// @Router       /api/v1/api-keys [post]
func CreateAPIKey(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		var req entity.CreateAPIKeyReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		resp, err := h.APIKey.Create(c.Request.Context(), userID, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, resp)
	}
}

// ListAPIKeys godoc
// @Summary      List API keys for the current user
// @Tags         api-keys
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  object{data=[]entity.APIKey}
// @Router       /api/v1/api-keys [get]
func ListAPIKeys(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		keys, err := h.APIKey.List(c.Request.Context(), userID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, keys)
	}
}

// RevokeAPIKey godoc
// @Summary      Revoke an API key
// @Tags         api-keys
// @Security     BearerAuth
// @Param        id  path  string  true  "API key ID"
// @Success      204
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/api-keys/{id} [delete]
func RevokeAPIKey(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		id := c.Param("id")
		if err := h.APIKey.Revoke(c.Request.Context(), id, userID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}
