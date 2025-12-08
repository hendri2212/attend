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
	query := h.db.Preload("School").Where("school_id = ?", schoolID).Order("name asc")

	query.Find(&classes)

	type ClassResponse struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	var response []ClassResponse
	for _, class := range classes {
		response = append(response, ClassResponse{
			ID:   class.ID,
			Name: class.Name,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *ClassHandler) CreateClass(c *gin.Context) {
	schoolIDInterface, exists := c.Get("school_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school_id not found in context"})
		return
	}
	schoolID := schoolIDInterface.(uint)

	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newClass := models.Class{
		Name:     input.Name,
		SchoolID: schoolID,
	}

	if err := h.db.Create(&newClass).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create class"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":   newClass.ID,
		"name": newClass.Name,
	})
}

func (h *ClassHandler) GetClassByID(c *gin.Context) {
	classID := c.Param("id")

	var class models.Class
	if err := h.db.Preload("School").First(&class, classID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":   class.ID,
		"name": class.Name,
	})
}

func (h *ClassHandler) UpdateClass(c *gin.Context) {
	classID := c.Param("id")

	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var class models.Class
	if err := h.db.First(&class, classID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
		return
	}

	class.Name = input.Name

	if err := h.db.Save(&class).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update class"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":   class.ID,
		"name": class.Name,
	})
}

func (h *ClassHandler) DeleteClass(c *gin.Context) {
	classID := c.Param("id")

	if err := h.db.Delete(&models.Class{}, classID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete class"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Class deleted successfully"})
}
