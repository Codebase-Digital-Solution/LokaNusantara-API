package router

import (
	"api/config"
	"api/handler"
	"api/middleware"
	"api/repository"
	"api/usecase"

	"time"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, ngrok-skip-browser-warning")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	locationRepo := repository.NewLocationRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	eventRepo := repository.NewEventRepository(db)
	trailRepo := repository.NewTrailRepository(db)

	// Initialize usecases
	authUsecase := usecase.NewAuthUsecase(userRepo)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	locationUsecase := usecase.NewLocationUsecase(locationRepo, categoryRepo)
	reviewUsecase := usecase.NewReviewUsecase(reviewRepo, userRepo, locationRepo)
	eventUsecase := usecase.NewEventUsecase(eventRepo, locationRepo)
	trailUsecase := usecase.NewTrailUsecase(trailRepo, locationRepo)

	// Tanya Loka is enabled only when an Anthropic API key is configured.
	var claude *anthropic.Client
	if cfg.AnthropicAPIKey != "" {
		client := anthropic.NewClient(option.WithAPIKey(cfg.AnthropicAPIKey))
		claude = &client
	}
	assistantUsecase := usecase.NewAssistantUsecase(db, claude)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUsecase, cfg)
	categoryHandler := handler.NewCategoryHandler(categoryUsecase)
	locationHandler := handler.NewLocationHandler(locationUsecase)
	reviewHandler := handler.NewReviewHandler(reviewUsecase)
	eventHandler := handler.NewEventHandler(eventUsecase)
	trailHandler := handler.NewTrailHandler(trailUsecase)
	assistantHandler := handler.NewAssistantHandler(assistantUsecase)
	adminHandler := handler.NewAdminHandler(db)

	// Public routes
	public := router.Group("/api/v1")
	{
		// Auth routes
		auth := public.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Public location routes
		locations := public.Group("/locations")
		{
			locations.GET("", locationHandler.GetAll)
			locations.GET("/search", locationHandler.Search)
			locations.GET("/nearby", locationHandler.GetNearby)
			locations.GET("/:id", locationHandler.GetByID)
			locations.GET("/:id/reviews", reviewHandler.GetLocationReviews)
			locations.GET("/:id/reviews/rating", reviewHandler.GetAverageRating)
		}

		// Public category routes
		categories := public.Group("/categories")
		{
			categories.GET("", categoryHandler.GetAll)
			categories.GET("/:id", categoryHandler.GetByID)
			categories.GET("/:id/locations", locationHandler.GetByCategory)
		}

		// Public cultural agenda
		public.GET("/events", eventHandler.GetAll)

		// Curated themed routes
		public.GET("/trails", trailHandler.GetAll)

		// Tanya Loka AI guide, rate limited per IP to control cost
		public.POST("/assistant/ask", middleware.RateLimit(20, 10*time.Minute), assistantHandler.Ask)
	}

	// Routes for any signed-in user (the mobile app)
	protected := router.Group("/api/v1")
	protected.Use(middleware.JWTMiddleware(cfg))
	{
		// Reviews: owners can edit or delete their own, admins any
		reviews := protected.Group("/reviews")
		{
			reviews.POST("", reviewHandler.Create)
			reviews.GET("/:id", reviewHandler.GetByID)
			reviews.PUT("/:id", reviewHandler.Update)
			reviews.DELETE("/:id", reviewHandler.Delete)
		}

		userReviews := protected.Group("/users/:userId/reviews")
		{
			userReviews.GET("", reviewHandler.GetUserReviews)
		}
	}

	// Admin-only routes: every change to shared data goes through here
	admin := router.Group("/api/v1")
	admin.Use(middleware.JWTMiddleware(cfg), middleware.RequireAdmin())
	{
		locations := admin.Group("/locations")
		{
			locations.POST("", locationHandler.Create)
			locations.PUT("/:id", locationHandler.Update)
			locations.DELETE("/:id", locationHandler.Delete)
		}

		categories := admin.Group("/categories")
		{
			categories.POST("", categoryHandler.Create)
			categories.PUT("/:id", categoryHandler.Update)
			categories.DELETE("/:id", categoryHandler.Delete)
		}

		events := admin.Group("/events")
		{
			events.POST("", eventHandler.Create)
			events.PUT("/:id", eventHandler.Update)
			events.DELETE("/:id", eventHandler.Delete)
		}

		trails := admin.Group("/trails")
		{
			trails.POST("", trailHandler.Create)
			trails.PUT("/:id", trailHandler.Update)
			trails.DELETE("/:id", trailHandler.Delete)
		}

		panel := admin.Group("/admin")
		{
			panel.GET("/stats", adminHandler.Stats)
			panel.GET("/users", adminHandler.Users)
			panel.GET("/events", eventHandler.AdminGetAll)
			panel.GET("/reviews", reviewHandler.AdminGetAll)
		}
	}

	return router
}
