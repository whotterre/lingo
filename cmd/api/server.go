package main

import (
	"lingo/internal/handlers"
	"lingo/pkg/auth/tokengen"
	"lingo/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	db "lingo/internal/db/sqlc"
)

type Server struct {
	router *gin.Engine
}

func dummy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello",
	})
}

func NewServer(pool *pgxpool.Pool, config utils.Config) *Server {
	router := gin.Default()

	server := &Server{}
	router.Use(gin.Logger())

	// Initialize handlers
	sqlStore := db.NewSQLStore(pool)
	newTok, err := tokengen.NewPasetoMaker(config.PasetoSecret)
	if err != nil {
		log.Fatal("Couldn't create token maker", err)
	}
	adminHandler := handlers.NewAdminHandler(sqlStore.(*db.SQLStore), newTok)

	public := router.Group("/v1/lingo")

	// Authentication routes
	public.POST("/auth/learner/signup", dummy)
	public.POST("/auth/learner/login", dummy)
	public.POST("/auth/learner/refresh", dummy)
	public.POST("/auth/admin/signup", adminHandler.RegisterAdmin)
	public.POST("/auth/admin/login", adminHandler.LoginAdmin)
	public.POST("/auth/admin/refresh", dummy)

	// Admin routes
	public.PUT("/admin/details/:adminId", adminHandler.UpdateAdminDetails)
	public.PUT("/admin/password/:adminId", adminHandler.UpdateAdminPassword)

	// Language routes
	public.POST("/admin/language/create", adminHandler.CreateNewLanguage)
	public.PUT("/admin/language/:languageId", adminHandler.UpdateLanguageById)
	public.DELETE("/admin/language/:languageId", adminHandler.DeleteLanguage)
	public.GET("/admin/lesson/languages/all", adminHandler.GetAvailableLanguages)

	// Course routes
	public.POST("/admin/course/create/:langId", adminHandler.CreateNewCourse)
	public.PUT("/admin/course/:courseId", adminHandler.UpdateCourseById)
	public.DELETE("/admin/course/:courseId", adminHandler.DeleteCourse)
	public.GET("/admin/lesson/courses/all", adminHandler.GetAllCourses)

	// Lesson routes
	public.POST("/admin/lesson/create/:courseId", adminHandler.CreateNewLesson)
	public.PUT("/admin/lesson/:lessonId", adminHandler.UpdateLessonById)
	public.DELETE("/admin/lesson/:lessonId", adminHandler.DeleteLesson)
	public.GET("/admin/lesson/lessons/all", adminHandler.GetAllLessons)
	public.GET("/admin/lesson/lessons/by-course/:courseId", adminHandler.GetLessonsByCourseId)

	// Exercise routes
	public.POST("/admin/exercise/create", adminHandler.CreateNewExercise)
	public.PUT("/admin/exercise/:exerciseId", adminHandler.UpdateExerciseById)
	public.DELETE("/admin/exercise/:exerciseId", adminHandler.DeleteExercise)
	public.GET("/admin/exercise/:exerciseId", adminHandler.GetExerciseById)
	public.GET("/admin/exercise/exercises/all", adminHandler.GetAllExercises)
	public.GET("/admin/exercise/exercises/by-lesson/:lessonId", adminHandler.GetExercisesByLessonId)

	// User routes
	public.GET("/users/me", dummy)
	public.PUT("/users/me", dummy)
	public.GET("/users/id", dummy)

	server.router = router
	return server
}
