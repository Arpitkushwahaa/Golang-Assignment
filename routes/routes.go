package routes

import (
	"stocky-backend/controllers"
	"stocky-backend/services"
	"stocky-backend/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures all routes and middleware
func SetupRouter() *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(utils.ErrorHandler())
	router.Use(utils.LoggingMiddleware())

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(config))

	// Initialize services
	priceService := services.NewPriceService()
	ledgerService := services.NewLedgerService()
	rewardService := services.NewRewardService(priceService, ledgerService)

	// Initialize controllers
	rewardController := controllers.NewRewardController(rewardService)

	// API routes
	api := router.Group("/api")
	{
		// Health check
		api.GET("/health", rewardController.HealthCheck)

		// Reward endpoints
		api.POST("/reward", rewardController.CreateReward)
		api.GET("/today-stocks/:userId", rewardController.GetTodayStocks)
		api.GET("/historical-inr/:userId", rewardController.GetHistoricalINR)
		api.GET("/stats/:userId", rewardController.GetStats)
		api.GET("/portfolio/:userId", rewardController.GetPortfolio)
	}

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Stocky Backend API",
			"version": "1.0.0",
			"endpoints": map[string]string{
				"POST /api/reward":                    "Create a new reward",
				"GET  /api/today-stocks/:userId":      "Get today's stock rewards",
				"GET  /api/historical-inr/:userId":    "Get historical INR valuations",
				"GET  /api/stats/:userId":             "Get user statistics",
				"GET  /api/portfolio/:userId":         "Get user portfolio",
				"GET  /api/health":                    "Health check",
			},
		})
	})

	return router
}
