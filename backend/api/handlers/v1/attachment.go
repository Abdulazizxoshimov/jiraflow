package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
)

// UploadAttachment godoc
// @Summary      Upload attachment
// @Tags         attachments
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        parent_type  path      string  true  "Parent type (issues|pages)"
// @Param        parent_id    path      string  true  "Parent ID"
// @Param        file         formData  file    true  "File to upload"
// @Success      201  {object}  object{data=entity.Attachment}
// @Failure      400  {object}  object{code=string,message=string}
// @Router       /api/v1/{parent_type}/{parent_id}/attachments [post]
func UploadAttachment(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		parentType := c.Param("parent_type")
		parentID := c.Param("parent_id")
		uploaderID := c.GetString(middleware.CtxUserID)

		fh, err := c.FormFile("file")
		if err != nil {
			hs.BadRequest(c, "file is required")
			return
		}
		f, err := fh.Open()
		if err != nil {
			hs.Error(c, err)
			return
		}
		defer f.Close()

		attachment, err := h.Attachment.Upload(
			c.Request.Context(),
			parentType, parentID, uploaderID,
			fh.Filename, fh.Size, fh.Header.Get("Content-Type"),
			f,
		)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, attachment)
	}
}

// GetAttachment godoc
// @Summary      Get attachment by ID
// @Tags         attachments
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Attachment ID"
// @Success      200  {object}  object{data=entity.Attachment}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/attachments/{id} [get]
func GetAttachment(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		attachment, err := h.Attachment.GetByID(c.Request.Context(), id)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, attachment)
	}
}

// ListAttachments godoc
// @Summary      List attachments for a parent
// @Tags         attachments
// @Produce      json
// @Security     BearerAuth
// @Param        parent_type  path  string  true  "Parent type (issues|pages)"
// @Param        parent_id    path  string  true  "Parent ID"
// @Success      200  {object}  object{data=[]entity.Attachment}
// @Router       /api/v1/{parent_type}/{parent_id}/attachments [get]
func ListAttachments(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		parentType := c.Param("parent_type")
		parentID := c.Param("parent_id")
		attachments, err := h.Attachment.ListByParent(c.Request.Context(), parentType, parentID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, attachments)
	}
}

// DeleteAttachment godoc
// @Summary      Delete attachment
// @Tags         attachments
// @Security     BearerAuth
// @Param        id  path  string  true  "Attachment ID"
// @Success      204
// @Failure      403  {object}  object{code=string,message=string}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/attachments/{id} [delete]
func DeleteAttachment(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		if err := h.Attachment.Delete(c.Request.Context(), id, actorID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// GetAttachmentURL godoc
// @Summary      Get presigned download URL for attachment
// @Tags         attachments
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Attachment ID"
// @Success      200  {object}  object{data=object{url=string}}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/attachments/{id}/url [get]
func GetAttachmentURL(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		url, err := h.Attachment.PresignedURL(c.Request.Context(), id)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, gin.H{"url": url})
	}
}
