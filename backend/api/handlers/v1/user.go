package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// ListUsers godoc
// @Summary      List users
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        page  query  int     false  "Page number"
// @Param        limit query  int     false  "Page size"
// @Success      200  {object}  object{data=[]entity.User,total=int,page=int,limit=int,total_pages=int}
// @Failure      401  {object}  object{code=string,message=string}
// @Router       /api/v1/users [get]
func ListUsers(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter entity.UserFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		users, total, err := h.User.List(c.Request.Context(), &filter)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.List(c, users, total, filter.Page, filter.GetLimit())
	}
}

// GetUser godoc
// @Summary      Get user by ID
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  string  true  "User ID"
// @Success      200  {object}  object{data=entity.User}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/users/{id} [get]
func GetUser(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		user, err := h.User.GetByID(c.Request.Context(), id)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, user)
	}
}

// CreateUser godoc
// @Summary      Create user
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  entity.CreateUserReq  true  "User data"
// @Success      201  {object}  object{data=entity.User}
// @Failure      400  {object}  object{code=string,message=string}
// @Router       /api/v1/users [post]
func CreateUser(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.CreateUserReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		user, err := h.User.Create(c.Request.Context(), &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, user)
	}
}

// UpdateUser godoc
// @Summary      Update user
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string               true  "User ID"
// @Param        body  body  entity.UpdateUserReq  true  "Update data"
// @Success      200  {object}  object{data=entity.User}
// @Failure      400  {object}  object{code=string,message=string}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/users/{id} [put]
func UpdateUser(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req entity.UpdateUserReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		user, err := h.User.Update(c.Request.Context(), id, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, user)
	}
}

// ChangePassword godoc
// @Summary      Change password
// @Tags         users
// @Accept       json
// @Security     BearerAuth
// @Param        id    path  string                    true  "User ID"
// @Param        body  body  entity.ChangePasswordReq  true  "Passwords"
// @Success      204
// @Failure      400  {object}  object{code=string,message=string}
// @Router       /api/v1/users/{id}/password [put]
func ChangePassword(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		var req entity.ChangePasswordReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		if err := h.User.ChangePassword(c.Request.Context(), userID, &req); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// DeactivateUser godoc
// @Summary      Deactivate user
// @Tags         users
// @Security     BearerAuth
// @Param        id  path  string  true  "User ID"
// @Success      204
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/users/{id}/deactivate [post]
func DeactivateUser(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := h.User.Deactivate(c.Request.Context(), id); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// ActivateUser godoc
// @Summary      Activate user
// @Tags         users
// @Security     BearerAuth
// @Param        id  path  string  true  "User ID"
// @Success      204
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/users/{id}/activate [post]
func ActivateUser(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := h.User.Activate(c.Request.Context(), id); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}
