package handler

import (
	"api/entity"
	"api/localize"
	"api/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TrailHandler struct {
	trailUsecase usecase.TrailUsecase
	db           *gorm.DB
}

func NewTrailHandler(trailUsecase usecase.TrailUsecase, db *gorm.DB) *TrailHandler {
	return &TrailHandler{trailUsecase: trailUsecase, db: db}
}

// TrailRequest is the admin create/update payload. Stops are stored in the
// order they are sent.
type TrailRequest struct {
	Slug     string `json:"slug" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Subtitle string `json:"subtitle"`
	Story    string `json:"story"`
	Color    string `json:"color"`
	ImageURL string `json:"image_url"`
	Stops    []struct {
		LocationID   string            `json:"location_id" binding:"required"`
		Note         string            `json:"note"`
		Translations map[string]string `json:"translations"`
	} `json:"stops" binding:"required,min=1,dive"`
	Translations map[string]entity.TrailTranslation `json:"translations"`
}

func (r TrailRequest) translations() (map[string]entity.TrailTranslation, []map[string]string) {
	values := r.Translations
	if values == nil {
		values = map[string]entity.TrailTranslation{}
	}
	if _, ok := values[localize.DefaultLanguage]; !ok {
		values[localize.DefaultLanguage] = entity.TrailTranslation{Title: r.Title, Subtitle: r.Subtitle, Story: r.Story}
	}
	stops := make([]map[string]string, len(r.Stops))
	for i, stop := range r.Stops {
		stops[i] = stop.Translations
		if stops[i] == nil {
			stops[i] = map[string]string{}
		}
		if _, ok := stops[i][localize.DefaultLanguage]; !ok {
			stops[i][localize.DefaultLanguage] = stop.Note
		}
	}
	return values, stops
}

func (r TrailRequest) toInput() usecase.TrailInput {
	in := usecase.TrailInput{
		Slug:     r.Slug,
		Title:    r.Title,
		Subtitle: r.Subtitle,
		Story:    r.Story,
		Color:    r.Color,
		ImageURL: r.ImageURL,
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
	include := c.Query("include_translations") == "true"
	for i := range trails {
		localize.Trail(&trails[i], c.Query("lang"), include)
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
	translations, stopTranslations := req.translations()
	if err := localize.SaveTrailTranslations(h.db, trail.ID, translations, stopTranslations); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save translations"})
		return
	}
	trail, _ = h.trailUsecase.GetTrail(trail.ID)
	localize.Trail(trail, c.Query("lang"), c.Query("include_translations") == "true")
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
	translations, stopTranslations := req.translations()
	if err := localize.SaveTrailTranslations(h.db, trail.ID, translations, stopTranslations); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save translations"})
		return
	}
	trail, _ = h.trailUsecase.GetTrail(trail.ID)
	localize.Trail(trail, c.Query("lang"), c.Query("include_translations") == "true")
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
