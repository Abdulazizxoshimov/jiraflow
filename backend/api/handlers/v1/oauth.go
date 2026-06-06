package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
)

// GoogleLogin godoc
// @Summary      Redirect to Google OAuth2 login
// @Tags         oauth
// @Produce      json
// @Param        redirect_url  query  string  false  "URL to redirect after login"
// @Success      302
// @Router       /api/v1/auth/google [get]
func GoogleLogin(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		redirectURL := c.Query("redirect_url")
		authURL, err := h.OAuth.GenerateAuthURL(c.Request.Context(), redirectURL)
		if err != nil {
			hs.Error(c, err)
			return
		}
		c.Redirect(http.StatusFound, authURL)
	}
}

// GoogleCallback godoc
// @Summary      Handle Google OAuth2 callback
// @Tags         oauth
// @Produce      json
// @Param        state  query  string  true  "OAuth state"
// @Param        code   query  string  true  "Authorization code"
// @Success      200  {object}  object{data=entity.OAuthCallbackResp}
// @Failure      400,401  {object}  object{code=string,message=string}
// @Router       /api/v1/auth/google/callback [get]
func GoogleCallback(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		state := c.Query("state")
		code := c.Query("code")
		if state == "" || code == "" {
			hs.BadRequest(c, "missing state or code")
			return
		}
		resp, err := h.OAuth.HandleCallback(c.Request.Context(), state, code)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, resp)
	}
}

// GoogleLink godoc
// @Summary      Link Google account to authenticated user
// @Tags         oauth
// @Security     BearerAuth
// @Param        redirect_url  query  string  false  "Redirect URL"
// @Success      302
// @Router       /api/v1/auth/google/link [get]
func GoogleLink(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		redirectURL := c.Query("redirect_url")
		authURL, err := h.OAuth.GenerateAuthURL(c.Request.Context(), redirectURL)
		if err != nil {
			hs.Error(c, err)
			return
		}
		c.Redirect(http.StatusFound, authURL)
	}
}

// GoogleLinkCallback godoc
// @Summary      Handle link callback (authenticated)
// @Tags         oauth
// @Security     BearerAuth
// @Success      204
// @Router       /api/v1/auth/google/link/callback [get]
func GoogleLinkCallback(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		state := c.Query("state")
		code := c.Query("code")
		if state == "" || code == "" {
			hs.BadRequest(c, "missing state or code")
			return
		}
		if err := h.OAuth.LinkAccount(c.Request.Context(), userID, state, code); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// UnlinkOAuthProvider godoc
// @Summary      Unlink a social login provider
// @Tags         oauth
// @Security     BearerAuth
// @Param        provider  path  string  true  "Provider (google)"
// @Success      204
// @Router       /api/v1/auth/providers/{provider} [delete]
func UnlinkOAuthProvider(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		provider := c.Param("provider")
		if err := h.OAuth.UnlinkAccount(c.Request.Context(), userID, provider); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// ListLinkedProviders godoc
// @Summary      List OAuth providers linked to the current user
// @Tags         oauth
// @Security     BearerAuth
// @Success      200  {object}  object{data=[]entity.OAuthAccount}
// @Router       /api/v1/auth/providers [get]
func ListLinkedProviders(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		accounts, err := h.OAuth.ListLinkedAccounts(c.Request.Context(), userID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, accounts)
	}
}
