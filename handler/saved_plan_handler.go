package handler

import (
	"api/usecase"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const maxPlanBodyBytes = 128 << 10 // 128 KB

// SavedPlanHandler serves the signed-in user's saved visit plans. The user
// always comes from the JWT.
type SavedPlanHandler struct {
	plans usecase.SavedPlanUsecase
}

func NewSavedPlanHandler(plans usecase.SavedPlanUsecase) *SavedPlanHandler {
	return &SavedPlanHandler{plans: plans}
}

type SavePlanRequest struct {
	LocationID string           `json:"location_id" binding:"required"`
	Plan       usecase.TripPlan `json:"plan" binding:"required"`
}

func (h *SavedPlanHandler) List(c *gin.Context) {
	plans, err := h.plans.List(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": plans})
}

func (h *SavedPlanHandler) Save(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxPlanBodyBytes)
	var req SavePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data rencana tidak valid"})
		return
	}
	saved, err := h.plans.Save(c.GetString("userID"), req.LocationID, req.Plan)
	var inputErr *usecase.PlanInputError
	switch {
	case errors.As(err, &inputErr):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "rencana gagal disimpan"})
	default:
		c.JSON(http.StatusCreated, saved)
	}
}

func (h *SavedPlanHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id rencana tidak valid"})
		return
	}
	err = h.plans.Delete(c.GetString("userID"), uint(id))
	switch {
	case errors.Is(err, usecase.ErrSavedPlanNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "rencana gagal dihapus"})
	default:
		c.JSON(http.StatusOK, gin.H{"message": "rencana dihapus"})
	}
}
