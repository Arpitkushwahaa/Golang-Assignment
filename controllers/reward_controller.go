package controllers

import (
	"net/http"
	"stocky-backend/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

// RewardController handles reward-related API endpoints
type RewardController struct {
	rewardService *services.RewardService
}

// NewRewardController creates a new reward controller
func NewRewardController(rewardService *services.RewardService) *RewardController {
	return &RewardController{
		rewardService: rewardService,
	}
}

// CreateRewardRequest represents the request body for POST /reward
type CreateRewardRequest struct {
	UserID    int     `json:"userId" binding:"required"`
	Symbol    string  `json:"symbol" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required,gt=0"`
	Timestamp string  `json:"timestamp" binding:"required"`
}

// CreateReward handles POST /reward
func (c *RewardController) CreateReward(ctx *gin.Context) {
	var req CreateRewardRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("Invalid request payload")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request payload: " + err.Error(),
		})
		return
	}

	// Parse timestamp
	timestamp, err := time.Parse(time.RFC3339, req.Timestamp)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid timestamp format, use RFC3339",
		})
		return
	}

	// Convert quantity to decimal
	quantity := decimal.NewFromFloat(req.Quantity)

	// Create reward
	err = c.rewardService.CreateReward(req.UserID, req.Symbol, quantity, timestamp)
	if err != nil {
		logrus.WithError(err).Error("Failed to create reward")
		
		// Check if it's a duplicate error
		if err.Error() == "duplicate reward: identical reward already exists" {
			ctx.JSON(http.StatusConflict, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create reward: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// GetTodayStocks handles GET /today-stocks/:userId
func (c *RewardController) GetTodayStocks(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID",
		})
		return
	}

	rewards, err := c.rewardService.GetTodayRewards(userID)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch today's rewards")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch today's rewards",
		})
		return
	}

	// Transform rewards to response format
	rewardsList := make([]map[string]interface{}, 0)
	for _, reward := range rewards {
		rewardsList = append(rewardsList, map[string]interface{}{
			"symbol":    reward.StockSymbol,
			"quantity":  reward.Quantity,
			"timestamp": reward.Timestamp.Format(time.RFC3339),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"userId":  userID,
		"rewards": rewardsList,
	})
}

// GetHistoricalINR handles GET /historical-inr/:userId
func (c *RewardController) GetHistoricalINR(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID",
		})
		return
	}

	days, err := c.rewardService.GetHistoricalINR(userID)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch historical INR")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch historical INR",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"userId": userID,
		"days":   days,
	})
}

// GetStats handles GET /stats/:userId
func (c *RewardController) GetStats(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID",
		})
		return
	}

	stats, err := c.rewardService.GetStats(userID)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch stats")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch stats",
		})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

// GetPortfolio handles GET /portfolio/:userId (BONUS endpoint)
func (c *RewardController) GetPortfolio(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID",
		})
		return
	}

	portfolio, err := c.rewardService.GetPortfolio(userID)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch portfolio")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch portfolio",
		})
		return
	}

	ctx.JSON(http.StatusOK, portfolio)
}

// HealthCheck handles GET /health
func (c *RewardController) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}
