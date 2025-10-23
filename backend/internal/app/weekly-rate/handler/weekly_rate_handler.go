package handler

import (
	"net/http"

	"app/internal/app/weekly-rate/model"
	WeeklyRateService "app/internal/app/weekly-rate/service"

	"github.com/gin-gonic/gin"
)

type WeeklyRateHandler struct {
	service WeeklyRateService.WeeklyRateService
}

func NewWeeklyRateHandler(service WeeklyRateService.WeeklyRateService) *WeeklyRateHandler {
	return &WeeklyRateHandler{service: service}
}

// GetAll Weekly Rates
//
// @Summary      Get all weekly rates
// @Description  Retrieve a list of all weekly rates
// @Tags         WeeklyRates
// @Produce      json
// @Success      201  {array}   model.WeeklyRate
// @Router       /weekly-rates [get]
func (handler *WeeklyRateHandler) GetAll(c *gin.Context) {
	_ = model.WeeklyRate{}
	weeklyRates, err := handler.service.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, weeklyRates)
}
