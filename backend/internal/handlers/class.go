package handlers

import (
	"attend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ClassHandler struct {
	db *gorm.DB
}

func ClassesHandler(db *gorm.DB) *ClassHandler {
	db.AutoMigrate(&models.Class{})
	return &ClassHandler{db: db}
}

func (h *ClassHandler) GetClassBySchool(c *gin.Context) {
	schoolIDInterface, exists := c.Get("school_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school_id not found in context"})
		return
	}
	schoolID := schoolIDInterface.(uint)

	var classes []models.Class
	query := h.db.Preload("School").Where("school_id = ?", schoolID)

	query.Find(&classes)

	type SchoolMinimal struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	type ClassResponse struct {
		ID        uint          `json:"id"`
		Name      string        `json:"name"`
		SchoolID  uint          `json:"school_id"`
		School    SchoolMinimal `json:"school"`
		CreatedAt string        `json:"created_at"`
		UpdatedAt string        `json:"updated_at"`
	}

	var response []ClassResponse
	for _, class := range classes {
		response = append(response, ClassResponse{
			ID:       class.ID,
			Name:     class.Name,
			SchoolID: class.SchoolID,
			School: SchoolMinimal{
				ID:   class.School.ID,
				Name: class.School.Name,
			},
			CreatedAt: class.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: class.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	c.JSON(http.StatusOK, response)
}
