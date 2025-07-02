package handlers

import (
	"attend/internal/models"
	"net/http"

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

// func (h *LeaderHandler) CreateLeader(c *gin.Context) {
// 	var leader models.Leader
// 	if err := c.ShouldBindJSON(&leader); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	h.db.Create(&leader)

// 	c.JSON(http.StatusCreated, leader)
// }
