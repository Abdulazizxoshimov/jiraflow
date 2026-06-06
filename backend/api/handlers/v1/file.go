package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
)

// UploadFile godoc
// @Summary      Upload a standalone file to object storage
// @Tags         files
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        file  formData  file  true  "File to upload"
// @Success      201  {object}  object{data=entity.File}
// @Failure      400  {object}  object{code=string,message=string}
// @Router       /api/v1/files/upload [post]
func UploadFile(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		attachment, err := h.File.Upload(
			c.Request.Context(),
			uploaderID,
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

// GetFilePresignedURL godoc
// @Summary      Get presigned URL for a stored file
// @Tags         files
// @Produce      json
// @Security     BearerAuth
// @Param        path  query  string  true  "Storage path returned from upload"
// @Success      200  {object}  object{data=object{url=string}}
// @Failure      400  {object}  object{code=string,message=string}
// @Router       /api/v1/files/presign [get]
func GetFilePresignedURL(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		storagePath := c.Query("path")
		if storagePath == "" {
			hs.BadRequest(c, "path is required")
			return
		}
		url, err := h.File.GetPresignedURL(c.Request.Context(), storagePath)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, gin.H{"url": url})
	}
}
