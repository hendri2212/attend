package seeds

import (
	"attend/internal/models"
	"time"

	"gorm.io/gorm"
)

func SeedSchools(db *gorm.DB) {
	var count int64
	db.Model(&models.School{}).Count(&count)
	if count == 0 {
		layout := "2006-01-02 15:04:05"
		t, _ := time.Parse(layout, "2025-05-30 11:29:52")

		school := models.School{
			ID:      1,
			NPSN:    "30303355",
			Name:    "SMPN 1 Kotabaru",
			Address: "Kotabaru",
			Logo:    "smpn1kotabaru.jpg",
			Email:   "administrator@smpn1kotabaru.sch.id",
			Telp:    "085746080544",
			Message: `{{.SchoolName}}
---------------------
Nama Lengkap: {{.FullName}}
Kelas: {{if .ClassName}}{{.ClassName}}{{else}}-{{end}}
Hari: {{.Day}}
Jam: {{.Time}}
---------------------
Terima kasih, anak Anda telah melakukan absensi hari ini!

*SPONSOR*
*Momo Donat | Sweet and Delicous*
*Belanja Lengkap Hanya di Winmart*
`,
			CreatedAt: t,
			UpdatedAt: t,
		}

		db.Create(&school)
	}
}
