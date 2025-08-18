package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"

	"{{MCP_MODULE_NAME}}/internal/config"
	"{{MCP_MODULE_NAME}}/internal/database"
	"{{MCP_MODULE_NAME}}/internal/handlers"
	"{{MCP_MODULE_NAME}}/internal/middleware"
	"{{MCP_MODULE_NAME}}/internal/services"
	"{{MCP_MODULE_NAME}}/pkg/logger"
	"{{MCP_MODULE_NAME}}/pkg/metrics"
)

func main() {
	logger.Init()
	log := logger.GetLogger()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration", "error", err)
	}

	metrics.Init()

	// Initialize PostgreSQL
	db, err := database.NewConnection(cfg.Database.URL)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run database migrations", "error", err)
	}

	// Initialize ClickHouse for analytics
	clickhouseDB, err := database.NewClickHouseConnection(cfg.ClickHouse.URL)
	if err != nil {
		log.Fatal("Failed to connect to ClickHouse", "error", err)
	}
	defer clickhouseDB.Close()

	// Initialize Redis for caching
	redisClient, err := database.NewRedisClient(cfg.Redis.URL)
	if err != nil {
		log.Fatal("Failed to connect to Redis", "error", err)
	}
	defer redisClient.Close()

	// Initialize AI services for {{MCP_DESCRIPTION}}
	aiService1, err := services.NewAI{{AI_SERVICE_1}}Service(cfg.AI)
	if err != nil {
		log.Fatal("Failed to initialize AI {{AI_SERVICE_1}}", "error", err)
	}

	aiService2, err := services.NewAI{{AI_SERVICE_2}}Service(cfg.AI)
	if err != nil {
		log.Fatal("Failed to initialize AI {{AI_SERVICE_2}}", "error", err)
	}

	aiService3, err := services.NewAI{{AI_SERVICE_3}}Service(cfg.AI)
	if err != nil {
		log.Fatal("Failed to initialize AI {{AI_SERVICE_3}}", "error", err)
	}

	aiService4, err := services.NewAI{{AI_SERVICE_4}}Service(cfg.AI)
	if err != nil {
		log.Fatal("Failed to initialize AI {{AI_SERVICE_4}}", "error", err)
	}

	// Initialize core services
	coreService := services.New{{CORE_SERVICE}}Service(db, clickhouseDB, aiService1, cfg)
	analyticsService := services.New{{ANALYTICS_SERVICE}}Service(clickhouseDB, aiService2, cfg)
	optimizationService := services.New{{OPTIMIZATION_SERVICE}}Service(db, aiService3, cfg)
	reportingService := services.New{{REPORTING_SERVICE}}Service(clickhouseDB, cfg)

	// Initialize handlers
	coreHandler := handlers.New{{CORE_SERVICE}}Handler(coreService)
	analyticsHandler := handlers.New{{ANALYTICS_SERVICE}}Handler(analyticsService)
	optimizationHandler := handlers.New{{OPTIMIZATION_SERVICE}}Handler(optimizationService)
	reportingHandler := handlers.New{{REPORTING_SERVICE}}Handler(reportingService)

	// Setup Gin router
	if cfg.Environment != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(logger.GinMiddleware())
	router.Use(metrics.GinMiddleware())

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = cfg.Security.AllowedOrigins
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Tenant-ID"}
	router.Use(cors.New(corsConfig))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "{{MCP_NAME}}",
			"version":   "1.0.0",
			"timestamp": time.Now().UTC(),
		})
	})

	// API routes
	api := router.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
	api.Use(middleware.TenantMiddleware())
	api.Use(middleware.RateLimitMiddleware(redisClient, cfg.RateLimit))
	{
		// {{CORE_FEATURE}} Management
		core := api.Group("/{{CORE_ENDPOINT}}")
		{
			core.GET("", coreHandler.List{{CORE_ENTITY}})
			core.POST("", coreHandler.Create{{CORE_ENTITY}})
			core.GET("/:id", coreHandler.Get{{CORE_ENTITY}})
			core.PUT("/:id", coreHandler.Update{{CORE_ENTITY}})
			core.DELETE("/:id", coreHandler.Delete{{CORE_ENTITY}})
			core.POST("/ai-optimize", coreHandler.AIOptimize{{CORE_ENTITY}})
		}

		// Analytics & Insights
		analytics := api.Group("/analytics")
		{
			analytics.GET("/metrics", analyticsHandler.GetMetrics)
			analytics.GET("/insights", analyticsHandler.GetInsights)
			analytics.POST("/ai-analysis", analyticsHandler.AIAnalysis)
			analytics.GET("/trends", analyticsHandler.GetTrends)
			analytics.POST("/reports", analyticsHandler.GenerateReport)
		}

		// Optimization with AI
		optimization := api.Group("/optimization")
		{
			optimization.GET("/recommendations", optimizationHandler.GetRecommendations)
			optimization.POST("/ai-optimize", optimizationHandler.AIOptimize)
			optimization.GET("/performance", optimizationHandler.GetPerformance)
			optimization.POST("/apply", optimizationHandler.ApplyOptimizations)
		}

		// Integration Hub
		integrations := api.Group("/integrations")
		{
			integrations.GET("/available", func(c *gin.Context) {
				integrations := map[string]interface{}{
					"{{INTEGRATION_CATEGORY_1}}": []string{"{{INTEGRATION_1}}", "{{INTEGRATION_2}}", "{{INTEGRATION_3}}"},
					"{{INTEGRATION_CATEGORY_2}}": []string{"{{INTEGRATION_4}}", "{{INTEGRATION_5}}", "{{INTEGRATION_6}}"},
				}
				c.JSON(http.StatusOK, integrations)
			})
			integrations.POST("/connect/:platform", coreHandler.ConnectPlatform)
			integrations.GET("/connected", coreHandler.GetConnectedPlatforms)
		}
	}

	// Admin routes
	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
	admin.Use(middleware.AdminMiddleware())
	{
		admin.GET("/dashboard", func(c *gin.Context) {
			stats := map[string]interface{}{
				"total_{{CORE_ENTITY_PLURAL}}":   coreService.GetTotal{{CORE_ENTITY}}(),
				"ai_optimizations_today":        aiService1.GetOptimizationsToday(),
				"performance_score":             optimizationService.GetPerformanceScore(),
				"active_integrations":           coreService.GetActiveIntegrations(),
			}
			c.JSON(http.StatusOK, stats)
		})

		admin.POST("/optimize-all", func(c *gin.Context) {
			go optimizationService.OptimizeAll()
			c.JSON(http.StatusAccepted, gin.H{"message": "System-wide optimization initiated"})
		})
	}

	// Metrics server
	metricsRouter := gin.New()
	metricsRouter.GET("/metrics", gin.WrapH(promhttp.Handler()))

	metricsServer := &http.Server{
		Addr:    ":" + cfg.MetricsPort,
		Handler: metricsRouter,
	}

	go func() {
		log.Info("Starting metrics server", "port", cfg.MetricsPort)
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Metrics server failed", "error", err)
		}
	}()

	// Main HTTP server
	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	go func() {
		log.Info("Starting HTTP server", "port", cfg.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("HTTP server failed", "error", err)
		}
	}()

	// Background tasks
	cronScheduler := cron.New(cron.WithSeconds())

	// Regular tasks every 30 minutes
	cronScheduler.AddFunc("0 */30 * * * *", func() {
		coreService.UpdateMetrics()
	})

	// AI optimization hourly
	cronScheduler.AddFunc("0 0 * * * *", func() {
		aiService1.RunOptimizations()
		aiService2.AnalyzeData()
	})

	// Generate daily reports at 7 AM
	cronScheduler.AddFunc("0 0 7 * * *", func() {
		reportingService.GenerateDailyReports()
	})

	// Cleanup old data weekly
	cronScheduler.AddFunc("0 0 0 * * 0", func() {
		coreService.CleanupOldData()
	})

	cronScheduler.Start()
	defer cronScheduler.Stop()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error("HTTP server forced to shutdown", "error", err)
	}

	if err := metricsServer.Shutdown(ctx); err != nil {
		log.Error("Metrics server forced to shutdown", "error", err)
	}

	log.Info("Server shutdown complete")
}