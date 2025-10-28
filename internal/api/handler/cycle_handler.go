package handler

import (
	"net/http"

	"github.com/bowe99/phone-usage-service/internal/application/dtos"
	"github.com/bowe99/phone-usage-service/internal/application/service"
	"github.com/gin-gonic/gin"
)

type CycleHandler struct {
	cycleService *service.CycleService
}

func SetupCycleHandler(cycleService *service.CycleService) *CycleHandler {
	return &CycleHandler{
		cycleService: cycleService,
	}
}

// GetCycleHistory handles POST /api/v1/cycles/history
// @Summary Get cycle history for an MDN
// @Description Retrieve the complete billing cycle history for a given MDN (phone number)
// @Tags cycles
// @Accept json
// @Produce json
// @Param request body dto.GetCycleHistoryRequest true "User ID and MDN"
// @Success 200 {array} model.CycleResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Router /api/v1/cycles/history [post]
func (h *CycleHandler) GetCycleHistory(c *gin.Context) {
	var req dto.GetCycleHistoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	cycles, err := h.cycleService.GetCycleHistory(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cycles": cycles,
	})
}
