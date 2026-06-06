package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/validator"
)

// ValidateBody binds JSON and runs struct validation.
// Returns the parsed request and true on success; writes the error response and returns false on failure.
func ValidateBody[T any](c *gin.Context) (*T, bool) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		hs.BadRequest(c, err.Error())
		return nil, false
	}
	if err := validator.Validate(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":   "VALIDATION_ERROR",
			"errors": err,
		})
		return nil, false
	}
	return &req, true
}
