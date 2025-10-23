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
// @Router       /users/weekly-rates [get]
func (handler *WeeklyRateHandler) GetAll(c *gin.Context) {
	_ = model.WeeklyRate{}
	weeklyRates, err := handler.service.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, weeklyRates)
}

// Create Weekly Rate
//
// @Summary      Create a new weekly rate
// @Description  Create a new weekly rate with the provided details
// @Tags         WeeklyRates
// @Accept       json
// @Produce      json
// @Param        weeklyRate  body      model.CreateWeeklyRate  true  "Weekly Rate Data"
// @Success      201 		"Weekly rate created successfully"
// @Router       /users/weekly-rates/create [post]
func (handler *WeeklyRateHandler) Create(c *gin.Context) {
	var newWeeklyRate model.CreateWeeklyRate
	if err := c.ShouldBindJSON(&newWeeklyRate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := handler.service.Create(newWeeklyRate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Weekly rate created successfully"})
}

// Update Weekly Rate
//
// @Summary      Update an existing weekly rate
// @Description  Update the details of an existing weekly rate
// @Tags         WeeklyRates
// @Accept       json
// @Produce      json
// @Param        uuid        path      string                  true  "Weekly Rate UUID"
// @Param        weeklyRate  body      model.UpdateWeeklyRate  true  "Weekly Rate Data"
// @Success      200 		"Weekly rate updated successfully"
// @Router       /users/weekly-rates/{uuid}/update [put]
func (handler *WeeklyRateHandler) Update(c *gin.Context) {
	uuid := c.Param("uuid")
	var updatedWeeklyRate model.UpdateWeeklyRate
	if err := c.ShouldBindJSON(&updatedWeeklyRate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := handler.service.Update(uuid, updatedWeeklyRate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Weekly rate updated successfully"})
}

// Delete Weekly Rate
//
// @Summary      Delete a weekly rate
// @Description  Delete a weekly rate by its UUID
// @Tags         WeeklyRates
// @Param        uuid  path      string  true  "Weekly Rate UUID"
// @Success      200 		"Weekly rate deleted successfully"
// @Router       /users/weekly-rates/{uuid}/delete [delete]
func (handler *WeeklyRateHandler) Delete(c *gin.Context) {
	uuid := c.Param("uuid")

	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID is required"})
		return
	}

	err := handler.service.Delete(uuid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Weekly rate deleted successfully"})
}
