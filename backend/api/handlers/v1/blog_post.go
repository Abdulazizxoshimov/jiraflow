package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

func CreateBlogPost(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.CreateBlogPostReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		bp, err := h.BlogPost.Create(c.Request.Context(), c.Param("id"), c.GetString("user_id"), &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		c.JSON(http.StatusCreated, bp)
	}
}

func ListBlogPosts(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		onlyPublished := c.Query("published") == "true"

		posts, total, err := h.BlogPost.List(c.Request.Context(), entity.ListBlogPostsFilter{
			SpaceID:       c.Param("id"),
			OnlyPublished: onlyPublished,
			Limit:         limit,
			Offset:        offset,
		})
		if err != nil {
			hs.Error(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"total": total, "items": posts})
	}
}

func GetBlogPost(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		bp, err := h.BlogPost.GetByID(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		c.JSON(http.StatusOK, bp)
	}
}

func UpdateBlogPost(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.UpdateBlogPostReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		bp, err := h.BlogPost.Update(c.Request.Context(), c.Param("id"), c.GetString("user_id"), &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		c.JSON(http.StatusOK, bp)
	}
}

func DeleteBlogPost(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.BlogPost.Delete(c.Request.Context(), c.Param("id"), c.GetString("user_id")); err != nil {
			hs.Error(c, err)
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func PublishBlogPost(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.BlogPost.Publish(c.Request.Context(), c.Param("id"), c.GetString("user_id")); err != nil {
			hs.Error(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"published": true})
	}
}

func UnpublishBlogPost(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.BlogPost.Unpublish(c.Request.Context(), c.Param("id"), c.GetString("user_id")); err != nil {
			hs.Error(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"published": false})
	}
}
