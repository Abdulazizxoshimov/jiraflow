package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// CreateCustomField godoc
// @Summary      Create custom field
// @Tags         custom-fields
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        project_id  path  string                      true  "Project ID"
// @Param        body        body  entity.CreateCustomFieldReq  true  "Custom field data"
// @Success      201  {object}  object{data=entity.CustomField}
// @Failure      400  {object}  object{code=string,message=string}
// @Router       /api/v1/projects/{project_id}/custom-fields [post]
func CreateCustomField(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		var req entity.CreateCustomFieldReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		cf := &entity.CustomField{
			Name:       req.Name,
			FieldKey:   req.FieldKey,
			FieldType:  req.FieldType,
			IsRequired: req.IsRequired,
			Options:    req.Options,
			Position:   req.Position,
		}
		field, err := h.CustomField.Create(c.Request.Context(), projectID, cf)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, field)
	}
}

// GetCustomField godoc
// @Summary      Get custom field by ID
// @Tags         custom-fields
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Custom field ID"
// @Success      200  {object}  object{data=entity.CustomField}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/custom-fields/{id} [get]
func GetCustomField(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		field, err := h.CustomField.GetByID(c.Request.Context(), id)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, field)
	}
}

// ListCustomFields godoc
// @Summary      List custom fields by project
// @Tags         custom-fields
// @Produce      json
// @Security     BearerAuth
// @Param        project_id  path  string  true  "Project ID"
// @Success      200  {object}  object{data=[]entity.CustomField}
// @Router       /api/v1/projects/{project_id}/custom-fields [get]
func ListCustomFields(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		fields, err := h.CustomField.ListByProject(c.Request.Context(), projectID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, fields)
	}
}

// UpdateCustomField godoc
// @Summary      Update custom field
// @Tags         custom-fields
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                      true  "Custom field ID"
// @Param        body  body  entity.UpdateCustomFieldReq  true  "Update data"
// @Success      200  {object}  object{data=entity.CustomField}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/custom-fields/{id} [put]
func UpdateCustomField(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req entity.UpdateCustomFieldReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		cf := &entity.CustomField{}
		if req.Name != nil {
			cf.Name = *req.Name
		}
		if req.IsRequired != nil {
			cf.IsRequired = *req.IsRequired
		}
		if req.Options != nil {
			cf.Options = req.Options
		}
		if req.Position != nil {
			cf.Position = *req.Position
		}
		field, err := h.CustomField.Update(c.Request.Context(), id, cf)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, field)
	}
}

// DeleteCustomField godoc
// @Summary      Delete custom field
// @Tags         custom-fields
// @Security     BearerAuth
// @Param        id  path  string  true  "Custom field ID"
// @Success      204
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/custom-fields/{id} [delete]
func DeleteCustomField(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := h.CustomField.Delete(c.Request.Context(), id); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// ReorderCustomFields godoc
// @Summary      Reorder custom fields
// @Tags         custom-fields
// @Accept       json
// @Security     BearerAuth
// @Param        project_id  path  string            true  "Project ID"
// @Param        body        body  object{id=int}    true  "Map of field_id to position"
// @Success      204
// @Router       /api/v1/projects/{project_id}/custom-fields/reorder [put]
func ReorderCustomFields(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		var positions map[string]int
		if err := c.ShouldBindJSON(&positions); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		if err := h.CustomField.Reorder(c.Request.Context(), projectID, positions); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}
