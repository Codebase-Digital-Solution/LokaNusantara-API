package handler

import (
	"api/entity"
	"api/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TrailHandler struct {
	trailUsecase usecase.TrailUsecase
}

func NewTrailHandler(trailUsecase usecase.TrailUsecase) *TrailHandler {
	return &TrailHandler{trailUsecase: trailUsecase}
}

// TrailRequest is the admin create/update payload. Stops are stored in the
// order they are sent.
type TrailRequest struct {
	Slug     string `json:"slug" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Subtitle string `json:"subtitle"`
	Story    string `json:"story"`
	Color    string `json:"color"`
	Stops    []struct {
		LocationID string `json:"location_id" binding:"required"`
		Note       string `json:"note"`
	} `json:"stops" binding:"required,min=1,dive"`
}

func (r TrailRequest) toInput() usecase.TrailInput {
	in := usecase.TrailInput{
		Slug:     r.Slug,
		Title:    r.Title,
		Subtitle: r.Subtitle,
		Story:    r.Story,
		Color:    r.Color,
	}
	for _, s := range r.Stops {
		in.Stops = append(in.Stops, usecase.TrailStopInput{LocationID: s.LocationID, Note: s.Note})
	}
	return in
}

func (h *TrailHandler) GetAll(c *gin.Context) {
	trails, err := h.trailUsecase.GetTrails()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if trails == nil {
		trails = make([]entity.Trail, 0)
	}
	c.JSON(http.StatusOK, gin.H{"data": trails})
}

func (h *TrailHandler) Create(c *gin.Context) {
	var req TrailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	trail, err := h.trailUsecase.CreateTrail(req.toInput())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, trail)
}

func (h *TrailHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trail id"})
		return
	}
	var req TrailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	trail, err := h.trailUsecase.UpdateTrail(uint(id), req.toInput())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, trail)
}

func (h *TrailHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trail id"})
		return
	}
	if err := h.trailUsecase.DeleteTrail(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "trail deleted successfully"})
}
