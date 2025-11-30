package handlers

import (
	"attend/internal/models"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (h *StudentHandler) GetStudentByRFID(c *gin.Context) {
	rfid := c.Param("rf_id")

	now := time.Now()
	if v, ok := h.inflight.Load(rfid); ok {
		if now.Sub(v.(time.Time)) < h.inflightDuration {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "duplicate request"})
			return
		}
	}
	h.inflight.Store(rfid, now)
	defer h.inflight.Delete(rfid)

	const sql = `
        SELECT
            s.user_id AS id, s.full_name, s.rf_id AS rfid,
            c.id   AS class_id,   c.name   AS class_name,
            p.id   AS parent_id,  p.whats_app AS parent_whatsapp
        FROM students s
        LEFT JOIN classes c ON s.class_id = c.id
        LEFT JOIN parents p ON s.parent_id = p.id
        WHERE s.rf_id = ?
        LIMIT 1
    `
	var res struct {
		ID             uint   `gorm:"column:id"`
		FullName       string `gorm:"column:full_name"`
		RFID           string `gorm:"column:rfid"`
		ClassID        uint   `gorm:"column:class_id"`
		ClassName      string `gorm:"column:class_name"`
		ParentID       uint   `gorm:"column:parent_id"`
		ParentWhatsApp string `gorm:"column:parent_whatsapp"`
	}
	err := h.db.Raw(sql, rfid).Scan(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	const duplicateGap = 5 * time.Minute
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if err := h.db.Transaction(func(tx *gorm.DB) error {
		var attendance models.Attendance
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&attendance, "student_id = ? AND date = ?", res.ID, today).
			Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			checkIn := now
			return tx.Create(&models.Attendance{
				StudentID: res.ID,
				Date:      today,
				CheckInAt: &checkIn,
				Method:    "rfid",
			}).Error
		}
		if err != nil {
			return err
		}

		if attendance.CheckOutAt != nil {
			return nil
		}
		if attendance.CheckInAt == nil {
			checkIn := now
			return tx.Model(&attendance).Updates(map[string]interface{}{
				"check_in_at": checkIn,
				"method":      "rfid",
			}).Error
		}
		if now.Sub(*attendance.CheckInAt) < duplicateGap {
			return nil
		}
		checkOut := now
		return tx.Model(&attendance).Updates(map[string]interface{}{
			"check_out_at": checkOut,
			"method":       "rfid",
		}).Error
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record attendance"})
		return
	}

	response := models.StudentResponse{
		ID:       res.ID,
		FullName: res.FullName,
		RFID:     res.RFID,
	}
	if res.ClassID != 0 {
		response.Class = &models.ClassResponse{
			ID:   res.ClassID,
			Name: res.ClassName,
		}
	}
	if res.ParentID != 0 && res.ParentWhatsApp != "" {
		response.Parent = &models.ParentResponse{
			ID:       res.ParentID,
			WhatsApp: res.ParentWhatsApp,
		}
	}
	c.JSON(http.StatusOK, response)

	student := models.Student{
		UserID:   res.ID,
		FullName: res.FullName,
		RFID:     res.RFID,
		Class: &models.Class{
			ID:   res.ClassID,
			Name: res.ClassName,
		},
		Parent: &models.Parent{
			ID:       res.ParentID,
			WhatsApp: res.ParentWhatsApp,
		},
	}

	select {
	case h.queue <- NotificationTask{Student: student, Retry: 0, NextDelay: 1 * time.Second}:
	default:
		// Queue full, drop silently or log if needed
	}
}

func (h *StudentHandler) GetAttendance(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 1000 { // Allow larger size for reports
		size = 20
	}
	offset := (page - 1) * size

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	classIDStr := c.Query("class_id")

	// If no date filter, use standard pagination
	if startDateStr == "" || endDateStr == "" {
		var attendances []models.Attendance
		query := h.db.
			Preload("Student").
			Preload("Student.Class").
			Order("date DESC, created_at DESC")

		if classIDStr != "" {
			query = query.Joins("JOIN students ON students.user_id = attendances.student_id").
				Where("students.class_id = ?", classIDStr)
		}

		if err := query.
			Limit(size).
			Offset(offset).
			Find(&attendances).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":      attendances,
			"page":      page,
			"page_size": size,
		})
		return
	}

	// If date filter exists, generate report including absent students
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
		return
	}

	// 1. Fetch Students
	var students []models.Student
	studentQuery := h.db.Preload("Class")
	if classIDStr != "" {
		studentQuery = studentQuery.Where("class_id = ?", classIDStr)
	}
	if err := studentQuery.Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch students"})
		return
	}

	// 2. Fetch Existing Attendance
	var attendances []models.Attendance
	attQuery := h.db.
		Preload("Student").
		Preload("Student.Class").
		Where("date >= ? AND date <= ?", startDate, endDate)

	if classIDStr != "" {
		attQuery = attQuery.Joins("JOIN students ON students.user_id = attendances.student_id").
			Where("students.class_id = ?", classIDStr)
	}

	if err := attQuery.Find(&attendances).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch attendance"})
		return
	}

	// Map: Date -> StudentID -> Attendance
	attMap := make(map[string]map[uint]models.Attendance)
	for _, att := range attendances {
		dateStr := att.Date.Format("2006-01-02")
		if _, ok := attMap[dateStr]; !ok {
			attMap[dateStr] = make(map[uint]models.Attendance)
		}
		attMap[dateStr][att.StudentID] = att
	}

	// 3. Generate Full Report
	var report []models.Attendance
	currentDate := startDate
	for !currentDate.After(endDate) {
		dateStr := currentDate.Format("2006-01-02")

		for _, student := range students {
			if att, ok := attMap[dateStr][student.UserID]; ok {
				report = append(report, att)
			} else {
				// Create Absent Record
				// Use a copy of student to avoid pointer issues if needed, but here it's fine
				s := student
				report = append(report, models.Attendance{
					Date:      currentDate,
					StudentID: s.UserID,
					Student:   &s,
					Method:    "", // Empty method implies absent/no record
				})
			}
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	// Pagination (In-Memory)
	total := len(report)
	start := offset
	end := offset + size
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	pagedData := report[start:end]

	c.JSON(http.StatusOK, gin.H{
		"data":      pagedData,
		"page":      page,
		"page_size": size,
		"total":     total,
	})
}
