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

type AssistantHandler struct {
	assistantUsecase usecase.AssistantUsecase
}

type AskRequest struct {
	Question string                     `json:"question" binding:"required"`
	History  []usecase.AssistantMessage `json:"history"`
}

func NewAssistantHandler(assistantUsecase usecase.AssistantUsecase) *AssistantHandler {
	return &AssistantHandler{assistantUsecase: assistantUsecase}
}

func (h *AssistantHandler) Ask(c *gin.Context) {
	var req AskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pertanyaan wajib diisi"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 90*time.Second)
	defer cancel()

	answer, err := h.assistantUsecase.Ask(ctx, req.Question, req.History)
	switch {
	case errors.Is(err, usecase.ErrAssistantDisabled):
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
	case err != nil && ctx.Err() != nil:
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Tanya Loka sedang lambat, coba lagi sebentar lagi"})
	case err != nil:
		log.Printf("[TanyaLoka] error: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Tanya Loka sedang tidak bisa menjawab, coba lagi nanti"})
	default:
		c.JSON(http.StatusOK, answer)
	}
}
