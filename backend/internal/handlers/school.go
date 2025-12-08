package handlers

import (
	"attend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SchoolHandler struct {
	db *gorm.DB
}

func SchoolsHandler(db *gorm.DB) *SchoolHandler {
	return &SchoolHandler{db: db}
}

func (h *SchoolHandler) GetSchools(c *gin.Context) {
	var schools []models.School
	if err := h.db.Find(&schools).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch schools"})
		return
	}
	c.JSON(http.StatusOK, schools)
}
