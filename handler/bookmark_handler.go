package handler

import (
	"api/entity"
	"api/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

// BookmarkHandler serves the signed-in user's saved places. The user always
// comes from the JWT, so nobody can read or change another user's list.
type BookmarkHandler struct {
	bookmarkUsecase usecase.BookmarkUsecase
}

func NewBookmarkHandler(bookmarkUsecase usecase.BookmarkUsecase) *BookmarkHandler {
	return &BookmarkHandler{bookmarkUsecase: bookmarkUsecase}
}

func (h *BookmarkHandler) GetAll(c *gin.Context) {
	locations, err := h.bookmarkUsecase.GetBookmarks(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if locations == nil {
		locations = make([]entity.Location, 0)
	}
	c.JSON(http.StatusOK, gin.H{"data": locations})
}

func (h *BookmarkHandler) Add(c *gin.Context) {
	if err := h.bookmarkUsecase.AddBookmark(c.GetString("userID"), c.Param("locationId")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"saved": true})
}

func (h *BookmarkHandler) Remove(c *gin.Context) {
	if err := h.bookmarkUsecase.RemoveBookmark(c.GetString("userID"), c.Param("locationId")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"saved": false})
}
