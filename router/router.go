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
	bookmarkRepo := repository.NewBookmarkRepository(db)

	// Initialize usecases
	authUsecase := usecase.NewAuthUsecase(userRepo)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	locationUsecase := usecase.NewLocationUsecase(locationRepo, categoryRepo)
	reviewUsecase := usecase.NewReviewUsecase(reviewRepo, userRepo, locationRepo)
	eventUsecase := usecase.NewEventUsecase(eventRepo, locationRepo)
	trailUsecase := usecase.NewTrailUsecase(trailRepo, locationRepo)
	bookmarkUsecase := usecase.NewBookmarkUsecase(bookmarkRepo, locationRepo)

	// OpenAI settings: admin-editable in the database, env as fallback.
	aiSettings := usecase.NewAISettingsUsecase(repository.NewSettingRepository(db), cfg.OpenAIAPIKey, cfg.OpenAIModel)

	// Tanya Loka needs an Anthropic or an OpenAI API key.
	var claude *anthropic.Client
	if cfg.AnthropicAPIKey != "" {
		client := anthropic.NewClient(option.WithAPIKey(cfg.AnthropicAPIKey))
		claude = &client
	}
	// Uses Claude when ANTHROPIC_API_KEY is set, otherwise OpenAI.
	assistantUsecase := usecase.NewAssistantUsecase(db, claude, aiSettings)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUsecase, cfg)
	categoryHandler := handler.NewCategoryHandler(categoryUsecase)
	locationHandler := handler.NewLocationHandler(locationUsecase)
	reviewHandler := handler.NewReviewHandler(reviewUsecase)
	eventHandler := handler.NewEventHandler(eventUsecase)
	trailHandler := handler.NewTrailHandler(trailUsecase)
	assistantHandler := handler.NewAssistantHandler(assistantUsecase)
	adminHandler := handler.NewAdminHandler(db)
	uploadHandler := handler.NewUploadHandler(cfg.UploadDir)
	adminSettingsHandler := handler.NewAdminSettingsHandler(aiSettings)
	plannerHandler := handler.NewPlannerHandler(
		usecase.NewPlannerUsecase(aiSettings, locationRepo, eventRepo),
	)
	bookmarkHandler := handler.NewBookmarkHandler(bookmarkUsecase)
	savedPlanHandler := handler.NewSavedPlanHandler(
		usecase.NewSavedPlanUsecase(repository.NewSavedPlanRepository(db), locationRepo),
	)

	// Uploaded images (location photos), public and cacheable.
	router.Static("/uploads", cfg.UploadDir)

	// Public routes
	public := router.Group("/api/v1")
	{
		// Auth routes
		auth := public.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/google", middleware.RateLimit(30, 10*time.Minute), authHandler.Google)
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

		// Visit planner (OpenAI), rate limited per IP to control cost
		public.POST("/planner", middleware.RateLimit(10, 10*time.Minute), plannerHandler.Plan)
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

		// Saved places of the signed-in user
		bookmarks := protected.Group("/bookmarks")
		{
			bookmarks.GET("", bookmarkHandler.GetAll)
			bookmarks.PUT("/:locationId", bookmarkHandler.Add)
			bookmarks.DELETE("/:locationId", bookmarkHandler.Remove)
		}

		// Saved visit plans of the signed-in user
		plans := protected.Group("/plans")
		{
			plans.GET("", savedPlanHandler.List)
			plans.POST("", savedPlanHandler.Save)
			plans.DELETE("/:id", savedPlanHandler.Delete)
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
			categories.PUT("/order", categoryHandler.Reorder)
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
			panel.POST("/uploads/location-image", uploadHandler.LocationImage)
			panel.POST("/uploads/trail-image", uploadHandler.TrailImage)
			panel.GET("/users", adminHandler.Users)
			panel.GET("/settings/ai", adminSettingsHandler.GetAI)
			panel.PUT("/settings/ai", adminSettingsHandler.UpdateAI)
			panel.POST("/settings/ai/test", adminSettingsHandler.TestAI)
			panel.GET("/events", eventHandler.AdminGetAll)
			panel.GET("/reviews", reviewHandler.AdminGetAll)
		}
	}

	return router
}
