package handler

import (
	"api/usecase"
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PlannerHandler struct {
	plannerUsecase usecase.PlannerUsecase
}

func NewPlannerHandler(plannerUsecase usecase.PlannerUsecase) *PlannerHandler {
	return &PlannerHandler{plannerUsecase: plannerUsecase}
}

// Plan generates a visit plan for one location within the visitor's budget.
func (h *PlannerHandler) Plan(c *gin.Context) {
	var req usecase.PlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data rencana tidak valid"})
		return
	}
	req.Language = c.GetHeader("Accept-Language")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 85*time.Second)
	defer cancel()

	plan, err := h.plannerUsecase.Plan(ctx, req)
	switch {
	case errors.Is(err, usecase.ErrPlannerDisabled):
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
	case err != nil && ctx.Err() != nil:
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Pembuatan rencana sedang lambat, coba lagi sebentar lagi"})
	case errors.As(err, new(*usecase.PlanInputError)):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case err != nil:
		log.Printf("[Planner] error: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Rencana belum bisa dibuat, coba lagi nanti"})
	default:
		c.JSON(http.StatusOK, plan)
	}
}
