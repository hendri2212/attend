package handlers

import (
	"attend/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"text/template"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var hariIndonesia = map[time.Weekday]string{
	time.Sunday:    "Minggu",
	time.Monday:    "Senin",
	time.Tuesday:   "Selasa",
	time.Wednesday: "Rabu",
	time.Thursday:  "Kamis",
	time.Friday:    "Jumat",
	time.Saturday:  "Sabtu",
}

func (h *StudentHandler) GetStudentByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	var student models.Student
	if err := h.db.
		Preload("Class").
		Preload("Parent").
		Preload("User").
		First(&student, "user_id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	response := models.StudentResponse{
		ID:         student.UserID,
		FullName:   student.FullName,
		WhatsApp:   student.WhatsApp,
		BirthPlace: student.BirthPlace,
		Born:       student.Born,
		RFID:       student.RFID,
		NISN:       student.NISN,
		Email:      student.User.Email,
	}
	if student.Class != nil {
		response.Class = &models.ClassResponse{
			ID:   student.Class.ID,
			Name: student.Class.Name,
		}
	}
	if student.Parent != nil {
		response.Parent = &models.ParentResponse{
			ID:       student.Parent.ID,
			FullName: student.Parent.FullName,
			WhatsApp: student.Parent.WhatsApp,
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *StudentHandler) GetStudents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	offset := (page - 1) * size

	var students []models.Student
	if err := h.db.
		Preload("Class").
		Preload("Parent").
		Limit(size).
		Offset(offset).
		Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	responses := make([]models.StudentResponse, 0, len(students))
	for _, s := range students {
		r := models.StudentResponse{
			ID:         s.UserID,
			FullName:   s.FullName,
			RFID:       s.RFID,
			WhatsApp:   s.WhatsApp,
			Born:       s.Born,
			BirthPlace: s.BirthPlace,
			NISN:       s.NISN,
		}
		if s.Class != nil {
			r.Class = &models.ClassResponse{
				ID:   s.Class.ID,
				Name: s.Class.Name,
			}
		}
		if s.Parent != nil {
			r.Parent = &models.ParentResponse{
				ID:       s.Parent.ID,
				FullName: s.Parent.FullName,
				WhatsApp: s.Parent.WhatsApp,
			}
		}
		responses = append(responses, r)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      responses,
		"page":      page,
		"page_size": size,
	})
}

func (h *StudentHandler) SaveStudent(c *gin.Context) {
	schoolIDVal, exists := c.Get("school_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "school_id not found in context"})
		return
	}
	schoolID, ok := schoolIDVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid school_id type"})
		return
	}

	type SaveStudentRequest struct {
		Email          string `json:"email" binding:"required,email"`
		NISN           string `json:"nisn" binding:"required"`
		FullName       string `json:"full_name" binding:"required"`
		RFID           string `json:"rfid"`
		WhatsApp       string `json:"whatsapp"`
		BirthPlace     string `json:"birth_place"`
		Born           string `json:"born" binding:"required"`
		ClassID        uint   `json:"class_id" binding:"required"`
		ParentName     string `json:"parent_name"`
		ParentWhatsApp string `json:"parent_whatsapp"`
	}

	var req SaveStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Convert WhatsApp if starts with '0' to '62...'
	if len(req.WhatsApp) > 0 && req.WhatsApp[0] == '0' {
		req.WhatsApp = "62" + req.WhatsApp[1:]
	}
	if len(req.ParentWhatsApp) > 0 && req.ParentWhatsApp[0] == '0' {
		req.ParentWhatsApp = "62" + req.ParentWhatsApp[1:]
	}

	bornTime, err := time.Parse("2006-01-02", req.Born)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format for born, use YYYY-MM-DD"})
		return
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	if err := h.db.Transaction(func(tx *gorm.DB) error {
		// 1. Handle Parent
		var parentID *uint
		if req.ParentWhatsApp != "" {
			var parent models.Parent
			// Check if parent exists by WhatsApp
			if err := tx.Where("whatsapp = ?", req.ParentWhatsApp).First(&parent).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Create new parent
					parent = models.Parent{
						FullName: req.ParentName,
						WhatsApp: req.ParentWhatsApp,
					}
					if err := tx.Create(&parent).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			}
			parentID = &parent.ID
		}

		user := models.User{
			Email:    req.Email,
			Password: string(hashedPwd),
			Role:     "student",
			SchoolID: schoolID,
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		student := models.Student{
			UserID:     user.ID,
			SchoolID:   schoolID,
			NISN:       req.NISN,
			FullName:   req.FullName,
			RFID:       req.RFID,
			ClassID:    req.ClassID,
			WhatsApp:   req.WhatsApp,
			BirthPlace: req.BirthPlace,
			Born:       bornTime,
			ParentID:   parentID,
		}
		if err := tx.Create(&student).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save student"})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	type UpdateStudentRequest struct {
		RFID       string `json:"rfid"`
		NISN       string `json:"nisn"`
		FullName   string `json:"full_name"`
		BirthPlace string `json:"birth_place"`
		Born       string `json:"born" binding:"required"`
		WhatsApp   string `json:"whatsapp"`
		ParentID   *uint  `json:"parent_id"`
		Email      string `json:"email"`
		ClassID    uint   `json:"class_id"`
	}

	var req UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.WhatsApp) > 0 && req.WhatsApp[0] == '0' {
		req.WhatsApp = "62" + req.WhatsApp[1:]
	}

	bornTime, err := time.Parse("2006-01-02", req.Born)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format for born, use YYYY-MM-DD"})
		return
	}

	if err := h.db.Transaction(func(tx *gorm.DB) error {
		var student models.Student
		if err := tx.Preload("User").First(&student, "user_id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("student not found")
			}
			return err
		}

		// Update User fields (Email)
		if req.Email != "" && student.User.Email != req.Email {
			student.User.Email = req.Email
			if err := tx.Save(student.User).Error; err != nil {
				return err
			}
		}

		// Update Student fields
		student.RFID = req.RFID
		student.NISN = req.NISN
		student.FullName = req.FullName
		student.BirthPlace = req.BirthPlace
		student.Born = bornTime
		student.WhatsApp = req.WhatsApp
		student.ParentID = req.ParentID

		if req.ClassID != 0 {
			student.ClassID = req.ClassID
		}

		if err := tx.Save(&student).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		if err.Error() == "student not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update student: " + err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	if err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.Student{}, "user_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.User{}, "id = ?", id).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete student"})
		return
	}

	c.Status(http.StatusOK)
}

var bulanIndonesia = map[time.Month]string{
	time.January:   "Januari",
	time.February:  "Februari",
	time.March:     "Maret",
	time.April:     "April",
	time.May:       "Mei",
	time.June:      "Juni",
	time.July:      "Juli",
	time.August:    "Agustus",
	time.September: "September",
	time.October:   "Oktober",
	time.November:  "November",
	time.December:  "Desember",
}

type NotificationTask struct {
	Student   models.Student
	Retry     int
	NextDelay time.Duration
}

type StudentHandler struct {
	db               *gorm.DB
	httpClient       *http.Client
	inflight         sync.Map
	inflightDuration time.Duration
	queue            chan NotificationTask
}

func StudentsHandler(db *gorm.DB) *StudentHandler {
	db.Logger = db.Logger.LogMode(logger.Warn)

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		},
	}

	h := &StudentHandler{
		db:               db,
		httpClient:       httpClient,
		inflightDuration: 1 * time.Second,
		queue:            make(chan NotificationTask, 100),
	}

	go h.worker()
	return h
}

func (h *StudentHandler) worker() {
	for task := range h.queue {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := h.sendWhatsAppNotification(ctx, task.Student)
		cancel()

		if err != nil && task.Retry < 3 {
			task.Retry++
			if task.NextDelay == 0 {
				task.NextDelay = 1 * time.Second
			} else {
				task.NextDelay *= 2
			}
			time.AfterFunc(task.NextDelay, func() { h.queue <- task })
		}
	}
}

func (h *StudentHandler) sendWhatsAppNotification(ctx context.Context, s models.Student) error {
	var school models.School
	if err := h.db.First(&school).Error; err != nil {
		return fmt.Errorf("failed to load school config: %w", err)
	}
	tmpl, err := template.New("whatsapp").Parse(school.Message)
	if err != nil {
		return fmt.Errorf("failed to parse message template: %w", err)
	}

	var buf bytes.Buffer
	now := time.Now()
	data := struct {
		SchoolName string
		FullName   string
		ClassName  string
		Day        string
		Time       string
	}{
		SchoolName: school.Name,
		FullName:   s.FullName,
		ClassName:  "",
		Day:        fmt.Sprintf("%s, %d %s %d", hariIndonesia[now.Weekday()], now.Day(), bulanIndonesia[now.Month()], now.Year()),
		Time:       now.Format("15:04:05"),
	}
	if s.Class != nil {
		data.ClassName = s.Class.Name
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	message := buf.String()

	payload := map[string]string{
		"to":      s.Parent.WhatsApp,
		"message": message,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://wabot.tukarjual.com/send", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("WhatsApp API returned status %d", resp.StatusCode)
	}
	return nil
}
