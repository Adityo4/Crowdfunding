package main

import (
	"net/http"
	"time"

	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/repository"
	"backend/internal/service"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

func main() {
	// Setup custom logger format
	log.SetTimeFormat(time.RFC3339)
	log.Info("Memulai inisialisasi server backend...")

	// Load Config
	cfg := config.LoadConfig()

	// Init PostgreSQL
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Gagal inisialisasi database PostgreSQL: %v", err)
	}

	// Init Redis
	rdb, err := database.InitRedis(cfg)
	if err != nil {
		log.Fatalf("Gagal inisialisasi Redis: %v", err)
	}
	_ = rdb

	// Init MinIO
	minioClient, err := database.InitMinio(cfg)
	if err != nil {
		log.Fatalf("Gagal inisialisasi MinIO: %v", err)
	}

	// Dependency Injection: Repositories
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	charityRepo := repository.NewCharityRepository(db)
	donationRepo := repository.NewDonationRepository(db)
	articleRepo := repository.NewArticleRepository(db)

	// Dependency Injection: Services
	authService := service.NewAuthService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	charityService := service.NewCharityService(charityRepo, minioClient)
	donationService := service.NewDonationService(donationRepo, charityRepo)
	articleService := service.NewArticleService(articleRepo)

	// Dependency Injection: Handlers
	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	charityHandler := handler.NewCharityHandler(charityService)
	donationHandler := handler.NewDonationHandler(donationService)
	articleHandler := handler.NewArticleHandler(articleService)

	// Set mode Gin
	gin.SetMode(gin.DebugMode)

	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Global Middleware (Recovery)
	r.Use(gin.Recovery())

	// Simple Ping Endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "healthy",
			"message":  "pong",
			"database": "connected",
			"redis":    "connected",
			"minio":    "connected",
		})
	})

	// API Routes Group
	v1 := r.Group("/v1")
	{
		// Auth Routes
		v1.POST("/auth/register", authHandler.Register)
		v1.POST("/auth/login", authHandler.Login)

		// User Routes
		v1.GET("/users/me", middleware.AuthMiddleware(), authHandler.Me)

		// Category Routes
		v1.GET("/categories", categoryHandler.GetCategories)

		// Charity / Campaign Routes
		v1.GET("/charities", charityHandler.GetCharities)
		v1.GET("/charities/:slug", charityHandler.GetCharityBySlug)
		v1.POST("/charities", middleware.AuthMiddleware(), charityHandler.CreateCharity)
		v1.PUT("/charities/:id/status", middleware.AuthMiddleware(), charityHandler.UpdateStatus)

		// Donation Routes
		v1.POST("/donations", middleware.OptionalAuthMiddleware(), donationHandler.CreateDonation)
		v1.POST("/donations/callback", donationHandler.ProcessCallback)

		// Article Routes
		v1.GET("/articles", articleHandler.GetArticles)
		v1.GET("/articles/:slug", articleHandler.GetArticleBySlug)
	}

	log.Infof("Server backend berhasil berjalan pada port %s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Gagal menjalankan server HTTP: %v", err)
	}
}
