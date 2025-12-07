package handlers

// import (
// 	"net/http"
// 	"attend/internal/models"

// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// )

// type PositionHandler struct {
// 	db *gorm.DB
// }

// func PositionsHandler(db *gorm.DB) *PositionHandler {
// 	db.AutoMigrate(&models.Position{})
// 	return &PositionHandler{db: db}
// }

// func (h *PositionHandler) GetPositions(c *gin.Context) {
// 	var positions []models.Position
// 	h.db.Find(&positions)

// 	c.JSON(http.StatusOK, positions)
// }

// func (h *PositionHandler) CreatePosition(c *gin.Context) {
// 	var position models.Position
// 	if err := c.ShouldBindJSON(&position); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	h.db.Create(&position)

// 	c.JSON(http.StatusCreated, position)
// }
