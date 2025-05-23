package main

import (
	"lingo/internal/handlers/auth"
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

func NewServer(pool *pgxpool.Pool) *Server {
	router := gin.Default()

	server := &Server{}
	router.Use(gin.Logger())
	// Initialize handlers here
	sqlStore := db.NewStore(pool)
	authHandler := handlers.NewAuthHandler(sqlStore.(*db.SQLStore))
	adminHandler := handlers.NewAdminHandler(sqlStore.(*db.SQLStore))
	public := router.Group("/v1/lingo")

	public.POST("/auth/learner/signup", authHandler.RegisterUser)
	public.POST("/auth/learner/login", authHandler.LoginUser)
	public.POST("/auth/learner/refresh", dummy)

	public.POST("/auth/admin/signup", adminHandler.SignUpAdmin)
	public.POST("/auth/admin/login", adminHandler.LoginAdmin)
	public.POST("/auth/admin/refresh", dummy)
	public.POST("/admin/language/create", adminHandler.CreateLanguage)
	public.POST("/admin/course/create/:langId", adminHandler.CreateCourse)
	


	public.GET("/users/me", dummy)
	public.PUT("/users/me", dummy)
	public.GET("/users/id", dummy)

	server.router = router
	return server
}
