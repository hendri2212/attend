package handlers

import (
	"attend/internal/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (h *StudentHandler) ImportStudents(c *gin.Context) {
	schoolIDVal, exists := c.Get("school_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "school_id not found in context"})
		return
	}
	schoolID := schoolIDVal.(uint)

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid excel file"})
		return
	}
	defer f.Close()

	// Get all rows in the first sheet.
	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty excel file"})
		return
	}
	rows, err := f.GetRows(sheetList[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get rows"})
		return
	}

	// Password default
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	successCount := 0
	errorsList := []string{}

	// Skip header row (index 0)
	for i, row := range rows {
		if i == 0 {
			// Header: email, nisn, full_name, born, class_id, rfid, whatsapp, birth_place, parent_name, parent_whatsapp
			continue
		}

		// Expected columns mapping (based on previous discussion)
		// 0: email (mandatory)
		// 1: nisn (mandatory)
		// 2: full_name (mandatory)
		// 3: born (mandatory YYYY-MM-DD)
		// 4: class_id (mandatory)
		// 5: rfid (optional)
		// 6: whatsapp (optional)
		// 7: birth_place (optional)
		// 8: parent_name (optional)
		// 9: parent_whatsapp (optional)

		if len(row) < 5 {
			errorsList = append(errorsList, fmt.Sprintf("Row %d: Incomplete data (min 5 columns required)", i+1))
			continue
		}

		email := strings.TrimSpace(row[0])
		nisn := strings.TrimSpace(row[1])
		fullName := strings.TrimSpace(row[2])
		bornStr := strings.TrimSpace(row[3])
		classIDStr := strings.TrimSpace(row[4])

		// Optional fields with safe access
		rfid := ""
		if len(row) > 5 {
			rfid = strings.TrimSpace(row[5])
		}
		whatsapp := ""
		if len(row) > 6 {
			whatsapp = strings.TrimSpace(row[6])
		}
		birthPlace := ""
		if len(row) > 7 {
			birthPlace = strings.TrimSpace(row[7])
		}
		parentName := ""
		if len(row) > 8 {
			parentName = strings.TrimSpace(row[8])
		}
		parentWhatsApp := ""
		if len(row) > 9 {
			parentWhatsApp = strings.TrimSpace(row[9])
		}

		if email == "" || nisn == "" || fullName == "" || bornStr == "" || classIDStr == "" {
			errorsList = append(errorsList, fmt.Sprintf("Row %d: Missing mandatory fields", i+1))
			continue
		}

		// Convert ClassID
		classID, err := strconv.Atoi(classIDStr)
		if err != nil {
			errorsList = append(errorsList, fmt.Sprintf("Row %d: Invalid Class ID", i+1))
			continue
		}

		// Convert Date
		bornTime, err := time.Parse("2006-01-02", bornStr)
		if err != nil {
			errorsList = append(errorsList, fmt.Sprintf("Row %d: Invalid date format (use YYYY-MM-DD)", i+1))
			continue
		}

		// Format WhatsApp
		if len(whatsapp) > 0 && whatsapp[0] == '0' {
			whatsapp = "62" + whatsapp[1:]
		}
		if len(parentWhatsApp) > 0 && parentWhatsApp[0] == '0' {
			parentWhatsApp = "62" + parentWhatsApp[1:]
		}

		// Start Transaction per Row
		err = h.db.Transaction(func(tx *gorm.DB) error {
			// 1. Handle Parent
			var parentID *uint
			if parentWhatsApp != "" {
				var parent models.Parent
				if err := tx.Where("whats_app = ?", parentWhatsApp).First(&parent).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						// Create New Parent
						parent = models.Parent{
							FullName: parentName,
							WhatsApp: parentWhatsApp,
						}
						if err := tx.Create(&parent).Error; err != nil {
							return fmt.Errorf("failed to create parent: %v", err)
						}
					} else {
						return err
					}
				}
				parentID = &parent.ID
			}

			// 2. Create User
			user := models.User{
				Email:    email,
				Password: string(hashedPwd),
				Role:     "student",
				SchoolID: schoolID,
			}
			if err := tx.Create(&user).Error; err != nil {
				return fmt.Errorf("failed to create user (email likely duplicate): %v", err)
			}

			// 3. Create Student
			student := models.Student{
				UserID:     user.ID,
				SchoolID:   schoolID,
				NISN:       nisn,
				FullName:   fullName,
				RFID:       rfid,
				ClassID:    uint(classID),
				WhatsApp:   whatsapp,
				BirthPlace: birthPlace,
				Born:       bornTime,
				ParentID:   parentID,
			}
			if err := tx.Create(&student).Error; err != nil {
				return fmt.Errorf("failed to create student: %v", err)
			}

			return nil
		})

		if err != nil {
			errorsList = append(errorsList, fmt.Sprintf("Row %d: %v", i+1, err))
		} else {
			successCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Import process completed",
		"success_count": successCount,
		"errors":        errorsList,
	})
}
