package handler

import (
	"api/usecase"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AdminSettingsHandler lets admins manage the OpenAI settings. The API key
// is write-only: responses only say whether one is set and its last 4
// characters.
type AdminSettingsHandler struct {
	ai usecase.AISettingsUsecase
}

func NewAdminSettingsHandler(ai usecase.AISettingsUsecase) *AdminSettingsHandler {
	return &AdminSettingsHandler{ai: ai}
}

func (h *AdminSettingsHandler) GetAI(c *gin.Context) {
	c.JSON(http.StatusOK, h.ai.View())
}

func (h *AdminSettingsHandler) UpdateAI(c *gin.Context) {
	var req usecase.AISettingsUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data pengaturan tidak valid"})
		return
	}
	view, err := h.ai.Update(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, view)
}

// TestAI checks the saved key and model against OpenAI.
func (h *AdminSettingsHandler) TestAI(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Second)
	defer cancel()
	if err := h.ai.Test(ctx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
