package seeds

import "gorm.io/gorm"

func SeedAll(db *gorm.DB) {
	SeedSchools(db)
	SeedClasses(db)
	SeedParents(db)
	SeedSuperAdmin(db)
	SeedStudents(db)
}
