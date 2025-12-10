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

	// Static routes (after CORS middleware)
	router.Static("/uploads", "./uploads")
	router.Static("/templates", "./templates")

	userHandler := handlers.UsersHandler(db)
	teachersHandler := handlers.TeachersHandler(db)
	parentsHandler := handlers.ParentsHandler(db)
	studentsHandler := handlers.StudentsHandler(db)
	classesHandler := handlers.ClassesHandler(db)
	schoolsHandler := handlers.SchoolsHandler(db)
	// budgetsHandler := handlers.BudgetsHandler(db)

	api := router.Group("/api")
	{
		api.POST("/login", userHandler.LoginUser)
		api.POST("/logout", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
		})
		api.GET("/student/:rf_id", studentsHandler.GetStudentByRFID)
		api.GET("/schools", schoolsHandler.GetSchools) // Public for login/register if needed, or moved to auth

		api.Use(middlewares.AuthMiddleware())

		api.GET("/me", userHandler.Me)
		api.GET("/users", userHandler.GetUsers)
		api.POST("/users", userHandler.CreateUser)
		api.GET("/users/:id", userHandler.GetUserByID)
		api.PUT("/users/:id", userHandler.UpdateUser)
		api.DELETE("/users/:id", userHandler.DeleteUser)

		api.GET("/teachers", teachersHandler.GetTeachers)
		api.POST("/teachers", teachersHandler.CreateTeacher)
		api.GET("/teachers/:id", teachersHandler.GetTeacherByID)
		api.PUT("/teachers/:id", teachersHandler.UpdateTeacher)
		// api.DELETE("/teachers/:id", userHandler.DeleteTeacher)

		api.GET("/students", studentsHandler.GetStudents)
		api.POST("/students", studentsHandler.SaveStudent)
		api.POST("/students/import", studentsHandler.ImportStudents)
		api.GET("/students/:id", studentsHandler.GetStudentByID)
		api.PUT("/students/:id", studentsHandler.UpdateStudent)
		api.DELETE("/students/:id", studentsHandler.DeleteStudent)

		api.GET("/classes", classesHandler.GetClassBySchool)
		api.POST("/classes", classesHandler.CreateClass)
		api.GET("/classes/:id", classesHandler.GetClassByID)
		api.PUT("/classes/:id", classesHandler.UpdateClass)
		api.DELETE("/classes/:id", classesHandler.DeleteClass)

		api.GET("/parents", parentsHandler.GetParents)
		api.POST("/parents", parentsHandler.CreateParent)
		api.GET("/parents/:id", parentsHandler.GetParentByID)
		api.PUT("/parents/:id", parentsHandler.UpdateParent)
		api.DELETE("/parents/:id", parentsHandler.DeleteParent)

		api.GET("/attendance", studentsHandler.GetAttendance)
	}
}
