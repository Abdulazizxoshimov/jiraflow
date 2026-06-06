package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// GetBoardSwimlanes godoc
// @Summary  Get board swimlanes
// @Tags     boards
// @Produce  json
// @Param    id        path   string  true  "Board ID"
// @Param    sprint_id query  string  false "Sprint ID filter"
// @Success  200       {object} entity.GetBoardSwimlanesResp
// @Router   /boards/{id}/swimlanes [get]
func GetBoardSwimlanes(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		boardID := c.Param("id")
		var sprintID *string
		if s := c.Query("sprint_id"); s != "" {
			sprintID = &s
		}
		resp, err := h.Board.GetBoardSwimlanes(c.Request.Context(), boardID, sprintID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// SetBoardSwimlaneType godoc
// @Summary  Set swimlane type for a board
// @Tags     boards
// @Accept   json
// @Produce  json
// @Param    id   path     string                    true  "Board ID"
// @Param    body body     entity.SetSwimlaneTypeReq true  "Swimlane type"
// @Success  200  {object} map[string]string
// @Router   /boards/{id}/swimlane-type [put]
func SetBoardSwimlaneType(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.SetSwimlaneTypeReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		if err := h.Board.SetSwimlaneType(c.Request.Context(), c.Param("id"), req.SwimlaneType); err != nil {
			hs.Error(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"swimlane_type": req.SwimlaneType})
	}
}
