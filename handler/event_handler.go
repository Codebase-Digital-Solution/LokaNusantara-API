package handler

import (
	"api/entity"
	"api/usecase"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const dateLayout = "2006-01-02"

type EventHandler struct {
	eventUsecase usecase.EventUsecase
}

// EventResponse sends dates as plain "YYYY-MM-DD" so clients never shift
// them across time zones.
type EventResponse struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Category    string           `json:"category"`
	StartDate   string           `json:"start_date"`
	EndDate     *string          `json:"end_date,omitempty"`
	IsEstimate  bool             `json:"is_estimate"`
	Region      string           `json:"region"`
	LocationID  *string          `json:"location_id,omitempty"`
	Location    *entity.Location `json:"location,omitempty"`
}

func NewEventHandler(eventUsecase usecase.EventUsecase) *EventHandler {
	return &EventHandler{eventUsecase: eventUsecase}
}

// GetAll lists events between ?from and ?to (YYYY-MM-DD). Defaults to the
// next 12 months starting today.
func (h *EventHandler) GetAll(c *gin.Context) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	from, err := parseDate(c.Query("from"), today)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'from' date, use YYYY-MM-DD"})
		return
	}
	to, err := parseDate(c.Query("to"), from.AddDate(1, 0, 0))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'to' date, use YYYY-MM-DD"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	events, err := h.eventUsecase.GetEvents(from, to, limit, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := make([]EventResponse, 0, len(events))
	for _, e := range events {
		data = append(data, toEventResponse(e))
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   data,
		"limit":  limit,
		"offset": offset,
	})
}

func toEventResponse(e entity.Event) EventResponse {
	item := EventResponse{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		Category:    e.Category,
		StartDate:   e.StartDate.Format(dateLayout),
		IsEstimate:  e.IsEstimate,
		Region:      e.Region,
		LocationID:  e.LocationID,
		Location:    e.Location,
	}
	if e.EndDate != nil {
		end := e.EndDate.Format(dateLayout)
		item.EndDate = &end
	}
	return item
}

// EventRequest is the admin create/update payload. Dates are YYYY-MM-DD.
type EventRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	StartDate   string  `json:"start_date" binding:"required"`
	EndDate     *string `json:"end_date"`
	IsEstimate  bool    `json:"is_estimate"`
	Region      string  `json:"region"`
	LocationID  *string `json:"location_id"`
}

func (r EventRequest) toInput() (usecase.EventInput, error) {
	start, err := time.ParseInLocation(dateLayout, r.StartDate, time.Local)
	if err != nil {
		return usecase.EventInput{}, errors.New("invalid start_date, use YYYY-MM-DD")
	}
	in := usecase.EventInput{
		Name:        r.Name,
		Description: r.Description,
		Category:    r.Category,
		StartDate:   start,
		IsEstimate:  r.IsEstimate,
		Region:      r.Region,
		LocationID:  r.LocationID,
	}
	if r.EndDate != nil && *r.EndDate != "" {
		end, err := time.ParseInLocation(dateLayout, *r.EndDate, time.Local)
		if err != nil {
			return usecase.EventInput{}, errors.New("invalid end_date, use YYYY-MM-DD")
		}
		in.EndDate = &end
	}
	return in, nil
}

// AdminGetAll lists every event regardless of date.
func (h *EventHandler) AdminGetAll(c *gin.Context) {
	events, err := h.eventUsecase.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	data := make([]EventResponse, 0, len(events))
	for _, e := range events {
		data = append(data, toEventResponse(e))
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *EventHandler) Create(c *gin.Context) {
	var req EventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	in, err := req.toInput()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	event, err := h.eventUsecase.CreateEvent(in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, toEventResponse(*event))
}

func (h *EventHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}
	var req EventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	in, err := req.toInput()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	event, err := h.eventUsecase.UpdateEvent(uint(id), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toEventResponse(*event))
}

func (h *EventHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}
	if err := h.eventUsecase.DeleteEvent(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "event deleted successfully"})
}

func parseDate(value string, fallback time.Time) (time.Time, error) {
	if value == "" {
		return fallback, nil
	}
	return time.ParseInLocation(dateLayout, value, time.Local)
}
