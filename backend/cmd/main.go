package main

import (
	"attend/internal/models"
	"attend/internal/routes"
	"attend/internal/seeds"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}
	return db
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}

	gin.SetMode(os.Getenv("GIN_MODE"))

	db := InitDB()

	// Auto migrate
	// Auto migrate
	if err := db.AutoMigrate(
		&models.School{},
		&models.Class{},
		&models.Parent{},
		&models.User{},
		&models.Student{},
		&models.Teacher{},
		&models.Attendance{},
	); err != nil {
		log.Fatal("failed to migrate database: ", err)
	}

	// Create Triggers
	if err := models.CreateTriggers(db); err != nil {
		log.Fatal("failed to create triggers: ", err)
	}

	// Seed data
	seeds.SeedAll(db)

	router := gin.Default()
	routes.SetupRoutes(router, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(router.Run(":" + port))
}
