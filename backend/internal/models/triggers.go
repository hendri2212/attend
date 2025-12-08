package models

import (
	"gorm.io/gorm"
)

// CreateTriggers creates database triggers to handle complex cascading deletions
// specifically ensuring Users are deleted when Students or Classes are deleted via DB.
func CreateTriggers(db *gorm.DB) error {
	// Trigger 1: When a Student is deleted, delete the associated User.
	// This covers "DELETE FROM students WHERE id = X"
	err := db.Exec(`
		DROP TRIGGER IF EXISTS after_student_delete;
	`).Error
	if err != nil {
		return err
	}

	err = db.Exec(`
		CREATE TRIGGER after_student_delete
		AFTER DELETE ON students
		FOR EACH ROW
		BEGIN
			DELETE FROM users WHERE id = OLD.user_id;
		END;
	`).Error
	if err != nil {
		return err
	}

	// Trigger 2: When a Class is deleted, delete associated Users (which cascades to Students).
	// This covers "DELETE FROM classes WHERE id = X"
	// Note: Standard Foreign Key CASCADE typically does NOT fire Triggers in MySQL.
	// So we must manually delete the Users BEFORE the Class is deleted.
	err = db.Exec(`
		DROP TRIGGER IF EXISTS before_class_delete;
	`).Error
	if err != nil {
		return err
	}

	// Use multi-table DELETE syntax compatible with MySQL
	err = db.Exec(`
		CREATE TRIGGER before_class_delete
		BEFORE DELETE ON classes
		FOR EACH ROW
		BEGIN
			DELETE users 
			FROM users 
			INNER JOIN students ON students.user_id = users.id 
			WHERE students.class_id = OLD.id;
		END;
	`).Error
	if err != nil {
		return err
	}

	return nil
}
