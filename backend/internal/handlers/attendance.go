package handlers

import (
	"attend/internal/models"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
