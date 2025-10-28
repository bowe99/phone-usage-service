package handler

import (
	"net/http"

	"github.com/bowe99/phone-usage-service/internal/application/dtos"
	"github.com/bowe99/phone-usage-service/internal/application/service"
	"github.com/gin-gonic/gin"
)

type DailyUsageHandler struct {
	dailyUsageService *service.DailyUsageService
}

func SetupDailyUsageHandler(usageService *service.DailyUsageService) *DailyUsageHandler {
	return &DailyUsageHandler{
		dailyUsageService: usageService,
	}
}

// GetCurrentCycleUsage handles POST /api/v1/usage/current-cycle
// @Summary Get current cycle daily usage
// @Description Retrieve daily usage data for the current billing cycle of a customer
// @Tags usage
// @Accept json
// @Produce json
// @Param request body dto.GetCurrentCycleUsageRequest true "User ID and MDN"
// @Success 200 {array} model.DailyUsageResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Router /api/v1/usage/current-cycle [post]
func (h *DailyUsageHandler) GetCurrentCycleUsage(c *gin.Context) {
	var req dto.GetCurrentCycleUsageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	usage, err := h.dailyUsageService.GetCurrentCycleUsage(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"usage": usage,
	})
}
