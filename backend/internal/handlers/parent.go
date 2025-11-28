package handlers

import (
	"attend/internal/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ParentHandler struct {
	db *gorm.DB
}

func ParentsHandler(db *gorm.DB) *ParentHandler {
	db.AutoMigrate(&models.Parent{})
	return &ParentHandler{db: db}
}

func (h *ParentHandler) GetParents(c *gin.Context) {
	var parents []models.Parent
	h.db.Find(&parents)

	c.JSON(http.StatusOK, parents)
}

func (h *ParentHandler) CreateParent(c *gin.Context) {
	var parent models.Parent
	if err := c.ShouldBindJSON(&parent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.HasPrefix(parent.WhatsApp, "0") {
		parent.WhatsApp = "62" + parent.WhatsApp[1:]
	}

	if err := h.db.Create(&parent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, parent)
}

func (h *ParentHandler) GetParentByID(c *gin.Context) {
	id := c.Param("id")
	var parent models.Parent
	if err := h.db.First(&parent, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Parent not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, parent)
}

func (h *ParentHandler) UpdateParent(c *gin.Context) {
	id := c.Param("id")
	var parent models.Parent
	if err := h.db.First(&parent, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Parent not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	var updatedData models.Parent
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.HasPrefix(updatedData.WhatsApp, "0") {
		updatedData.WhatsApp = "62" + updatedData.WhatsApp[1:]
	}

	if err := h.db.Model(&parent).Updates(updatedData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, parent)
}

func (h *ParentHandler) DeleteParent(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&models.Parent{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Parent deleted successfully"})
}
