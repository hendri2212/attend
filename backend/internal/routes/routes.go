package routes

import (
	"attend/internal/handlers"
	"attend/internal/middlewares"
	"net/http"

	"gorm.io/gorm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	router.Static("/uploads", "./uploads")

	var allowedOrigins []string
	if gin.Mode() == gin.ReleaseMode {
		allowedOrigins = []string{
			"https://attend.saijaan.com",
		}
	} else {
		allowedOrigins = []string{
			"http://localhost:5173",
		}
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		// AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNoContent)
	})

	userHandler := handlers.UsersHandler(db)
	// positionsHandler := handlers.PositionsHandler(db)
	parentsHandler := handlers.ParentsHandler(db)
	studentsHandler := handlers.StudentsHandler(db)
	classesHandler := handlers.ClassesHandler(db)
	// budgetsHandler := handlers.BudgetsHandler(db)

	api := router.Group("/api")
	{
		api.POST("/login", userHandler.LoginUser)
		api.POST("/logout", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
		})
		api.GET("/student/:rf_id", studentsHandler.GetStudentByRFID)

		api.Use(middlewares.AuthMiddleware())

		// api.GET("/me", userHandler.Me)
		// api.GET("/users", userHandler.GetUsers)
		// api.POST("/users", userHandler.CreateUser)
		// api.GET("/users/:id", userHandler.GetUserByID)
		// api.DELETE("/users/:id", userHandler.DeleteUser)

		api.GET("/students", studentsHandler.GetStudents)
		api.POST("/students", studentsHandler.SaveStudent)
		api.GET("/students/:id", studentsHandler.GetStudentByID)
		api.PUT("/students/:id", studentsHandler.UpdateStudent)

		api.GET("/classes", classesHandler.GetClassBySchool)
		api.GET("/parents", parentsHandler.GetParents)
	}
}
