package handlers

import (
	"attend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type TeacherHandler struct {
	db *gorm.DB
}

func TeachersHandler(db *gorm.DB) *TeacherHandler {
	return &TeacherHandler{db: db}
}

func (h *TeacherHandler) GetTeachers(c *gin.Context) {
	var teachers []models.Teacher
	h.db.Preload("User").Find(&teachers)

	c.JSON(http.StatusOK, teachers)
}

func (h *TeacherHandler) CreateTeacher(c *gin.Context) {
	// Request payload structure
	var req struct {
		SchoolID uint    `json:"school_id" binding:"required"`
		FullName string  `json:"full_name" binding:"required"`
		Nip      *string `json:"nip"`
		Email    string  `json:"email" binding:"required,email"`
		Password string  `json:"password"` // ignored, using hardcoded password
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start transaction
	tx := h.db.Begin()

	// Hash the default password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user with hashed password
	user := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     models.UserRoleTeacher,
		SchoolID: req.SchoolID,
		IsActive: true,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	// Create teacher with user_id from created user
	teacher := models.Teacher{
		UserID:   user.ID,
		FullName: req.FullName,
		Nip:      req.Nip,
	}

	if err := tx.Create(&teacher).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create teacher: " + err.Error()})
		return
	}

	// Commit transaction
	tx.Commit()

	// Load user relation for response
	h.db.Preload("User").First(&teacher, teacher.UserID)

	c.JSON(http.StatusCreated, teacher)
}

func (h *TeacherHandler) GetTeacherByID(c *gin.Context) {
	var teacher models.Teacher
	if err := h.db.Preload("User").First(&teacher, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher not found"})
		return
	}

	c.JSON(http.StatusOK, teacher)
}

func (h *TeacherHandler) UpdateTeacher(c *gin.Context) {
	var teacher models.Teacher
	if err := h.db.Preload("User").First(&teacher, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher not found"})
		return
	}

	// Request payload structure
	var req struct {
		SchoolID uint    `json:"school_id" binding:"required"`
		FullName string  `json:"full_name" binding:"required"`
		Nip      *string `json:"nip"`
		Email    string  `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update teacher
	teacher.FullName = req.FullName
	teacher.Nip = req.Nip
	teacher.User.SchoolID = req.SchoolID
	teacher.User.Email = req.Email

	// Save both teacher and user
	if err := h.db.Save(&teacher.User).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user: " + err.Error()})
		return
	}
	if err := h.db.Save(&teacher).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update teacher: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, teacher)
}
