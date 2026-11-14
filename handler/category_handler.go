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

type CategoryHandler struct {
	categoryUsecase usecase.CategoryUsecase
	db              *gorm.DB
}

type CategoryRequest struct {
	Name         string                                `json:"name"`
	Translations map[string]entity.CategoryTranslation `json:"translations"`
}

func NewCategoryHandler(categoryUsecase usecase.CategoryUsecase, db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{categoryUsecase: categoryUsecase, db: db}
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Name == "" {
		req.Name = req.Translations[localize.DefaultLanguage].Label
	}

	category, err := h.categoryUsecase.CreateCategory(req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	translations := req.Translations
	if translations == nil {
		translations = map[string]entity.CategoryTranslation{}
	}
	if _, ok := translations[localize.DefaultLanguage]; !ok {
		translations[localize.DefaultLanguage] = entity.CategoryTranslation{Label: req.Name}
	}
	if err := localize.SaveCategoryTranslations(h.db, category.ID, translations); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save translations"})
		return
	}
	category, _ = h.categoryUsecase.GetCategory(category.ID)
	localize.Category(category, c.Query("lang"), c.Query("include_translations") == "true")
	c.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	category, err := h.categoryUsecase.GetCategory(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	localize.Category(category, c.Query("lang"), c.Query("include_translations") == "true")
	c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name == "" {
		if old, getErr := h.categoryUsecase.GetCategory(uint(id)); getErr == nil {
			req.Name = old.Name
		}
	}
	category, err := h.categoryUsecase.UpdateCategory(uint(id), req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	translations := req.Translations
	if translations == nil {
		translations = map[string]entity.CategoryTranslation{}
	}
	if _, ok := translations[localize.DefaultLanguage]; !ok {
		translations[localize.DefaultLanguage] = entity.CategoryTranslation{Label: req.Name}
	}
	if err := localize.SaveCategoryTranslations(h.db, uint(id), translations); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save translations"})
		return
	}
	category, _ = h.categoryUsecase.GetCategory(uint(id))
	localize.Category(category, c.Query("lang"), c.Query("include_translations") == "true")
	c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	if err := h.categoryUsecase.DeleteCategory(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted successfully"})
}

type CategoryOrderRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}

// Reorder saves the display order of all categories (admin only).
func (h *CategoryHandler) Reorder(c *gin.Context) {
	var req CategoryOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "urutan kategori tidak valid"})
		return
	}
	if err := h.categoryUsecase.ReorderCategories(req.IDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "urutan disimpan"})
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")

	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)

	categories, err := h.categoryUsecase.GetAllCategories(limitInt, offsetInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range categories {
		localize.Category(&categories[i], c.Query("lang"), c.Query("include_translations") == "true")
	}
	if categories == nil {
		categories = make([]entity.Category, 0)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   categories,
		"limit":  limitInt,
		"offset": offsetInt,
	})
}
