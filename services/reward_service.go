package services

import (
	"errors"
	"fmt"
	"stocky-backend/db"
	"stocky-backend/models"
	"stocky-backend/utils"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// RewardService handles reward operations
type RewardService struct {
	priceService  *PriceService
	ledgerService *LedgerService
}

// NewRewardService creates a new reward service
func NewRewardService(priceService *PriceService, ledgerService *LedgerService) *RewardService {
	return &RewardService{
		priceService:  priceService,
		ledgerService: ledgerService,
	}
}

// CreateReward creates a new reward event with ledger entries
func (s *RewardService) CreateReward(userID int, symbol string, quantity decimal.Decimal, timestamp time.Time) error {
	// Validate inputs
	if err := utils.ValidateStockSymbol(symbol); err != nil {
		return fmt.Errorf("invalid stock symbol: %w", err)
	}

	if err := utils.ValidateQuantity(quantity); err != nil {
		return fmt.Errorf("invalid quantity: %w", err)
	}

	// Round quantity to 6 decimal places
	quantity = quantity.Round(6)

	// Check for duplicate reward (deduplication)
	var existingReward models.RewardEvent
	err := db.DB.Where("user_id = ? AND stock_symbol = ? AND quantity = ? AND timestamp = ?",
		userID, symbol, quantity, timestamp).First(&existingReward).Error

	if err == nil {
		logrus.WithFields(logrus.Fields{
			"userId":    userID,
			"symbol":    symbol,
			"quantity":  quantity,
			"timestamp": timestamp,
		}).Warn("Duplicate reward detected, rejecting")
		return fmt.Errorf("duplicate reward: identical reward already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check for duplicates: %w", err)
	}

	// Get current price for the stock
	pricePerShare, err := s.priceService.GetCurrentPrice(symbol)
	if err != nil {
		return fmt.Errorf("failed to get stock price: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"userId":       userID,
		"symbol":       symbol,
		"quantity":     quantity,
		"pricePerShare": pricePerShare,
		"timestamp":    timestamp,
	}).Info("Creating reward event")

	// Start transaction
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create reward event
	rewardEvent := models.RewardEvent{
		UserID:      userID,
		StockSymbol: symbol,
		Quantity:    quantity,
		Timestamp:   timestamp,
	}

	if err := tx.Create(&rewardEvent).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create reward event: %w", err)
	}

	// Create ledger entries
	if err := s.createLedgerEntriesInTx(tx, &rewardEvent, pricePerShare); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create ledger entries: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"rewardId": rewardEvent.ID,
		"userId":   userID,
		"symbol":   symbol,
		"quantity": quantity,
	}).Info("Reward created successfully")

	return nil
}

// createLedgerEntriesInTx creates ledger entries within a transaction
func (s *RewardService) createLedgerEntriesInTx(tx *gorm.DB, rewardEvent *models.RewardEvent, pricePerShare decimal.Decimal) error {
	quantity := rewardEvent.Quantity
	symbol := rewardEvent.StockSymbol
	timestamp := rewardEvent.Timestamp

	totalValue := utils.RoundINR(pricePerShare.Mul(quantity))
	_, _, _, totalFees := utils.CalculateFees(pricePerShare, quantity)

	entries := []models.LedgerEntry{
		{
			RewardEventID: rewardEvent.ID,
			EntryType:     models.EntryTypeStock,
			StockSymbol:   &symbol,
			Quantity:      quantity,
			AmountINR:     totalValue,
			Timestamp:     timestamp,
		},
		{
			RewardEventID: rewardEvent.ID,
			EntryType:     models.EntryTypeCash,
			StockSymbol:   &symbol,
			Quantity:      decimal.Zero,
			AmountINR:     totalValue.Neg(),
			Timestamp:     timestamp,
		},
		{
			RewardEventID: rewardEvent.ID,
			EntryType:     models.EntryTypeFee,
			StockSymbol:   &symbol,
			Quantity:      decimal.Zero,
			AmountINR:     totalFees.Neg(),
			Timestamp:     timestamp,
		},
	}

	return tx.Create(&entries).Error
}

// GetTodayRewards retrieves all reward events for a user for today
func (s *RewardService) GetTodayRewards(userID int) ([]models.RewardEvent, error) {
	now := utils.NowUTC()
	startOfDay := utils.StartOfDayUTC(now)
	endOfDay := utils.EndOfDayUTC(now)

	var rewards []models.RewardEvent
	err := db.DB.Where("user_id = ? AND timestamp >= ? AND timestamp <= ?",
		userID, startOfDay, endOfDay).
		Order("timestamp DESC").
		Find(&rewards).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch today's rewards: %w", err)
	}

	return rewards, nil
}

// GetHistoricalINR calculates INR valuation per past day (up to yesterday)
func (s *RewardService) GetHistoricalINR(userID int) ([]map[string]interface{}, error) {
	yesterday := utils.GetYesterday()
	
	// Get the earliest reward date for this user
	var firstReward models.RewardEvent
	err := db.DB.Where("user_id = ?", userID).
		Order("timestamp ASC").
		First(&firstReward).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []map[string]interface{}{}, nil
		}
		return nil, fmt.Errorf("failed to fetch first reward: %w", err)
	}

	startDate := utils.StartOfDayUTC(firstReward.Timestamp)
	
	// Generate list of dates from first reward to yesterday
	dates := utils.GetPastDates(startDate, yesterday)

	var result []map[string]interface{}

	for _, date := range dates {
		endOfDate := utils.EndOfDayUTC(date)
		
		// Get holdings up to this date
		holdings, err := s.ledgerService.GetUserStockHoldingsUpToDate(userID, endOfDate)
		if err != nil {
			logrus.Errorf("Failed to get holdings for date %s: %v", utils.GetDateString(date), err)
			continue
		}

		// Calculate total INR value for this date
		totalValue := decimal.Zero
		
		for symbol, qty := range holdings {
			price, err := s.priceService.GetPriceAtTime(symbol, endOfDate)
			if err != nil {
				logrus.Warnf("Failed to get price for %s at %s: %v", symbol, utils.GetDateString(date), err)
				continue
			}
			value := price.Mul(qty)
			totalValue = totalValue.Add(value)
		}

		result = append(result, map[string]interface{}{
			"date":     utils.GetDateString(date),
			"valueINR": utils.RoundINR(totalValue),
		})
	}

	return result, nil
}

// GetStats returns today's reward stats and current portfolio value
func (s *RewardService) GetStats(userID int) (map[string]interface{}, error) {
	// Get today's rewards grouped by stock
	todayRewards, err := s.GetTodayRewards(userID)
	if err != nil {
		return nil, err
	}

	// Group by stock symbol
	todayRewardsByStock := make(map[string]decimal.Decimal)
	for _, reward := range todayRewards {
		current, exists := todayRewardsByStock[reward.StockSymbol]
		if exists {
			todayRewardsByStock[reward.StockSymbol] = current.Add(reward.Quantity)
		} else {
			todayRewardsByStock[reward.StockSymbol] = reward.Quantity
		}
	}

	// Get total portfolio holdings
	holdings, err := s.ledgerService.GetUserStockHoldings(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get holdings: %w", err)
	}

	// Calculate current portfolio value
	totalValue := decimal.Zero
	for symbol, qty := range holdings {
		price, err := s.priceService.GetCurrentPrice(symbol)
		if err != nil {
			logrus.Warnf("Failed to get current price for %s: %v", symbol, err)
			continue
		}
		value := price.Mul(qty)
		totalValue = totalValue.Add(value)
	}

	return map[string]interface{}{
		"userId":            userID,
		"todayRewards":      todayRewardsByStock,
		"portfolioValueINR": utils.RoundINR(totalValue),
	}, nil
}

// GetPortfolio returns full holdings with current INR value
func (s *RewardService) GetPortfolio(userID int) (map[string]interface{}, error) {
	holdings, err := s.ledgerService.GetUserStockHoldings(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get holdings: %w", err)
	}

	var portfolioItems []map[string]interface{}
	totalValue := decimal.Zero

	for symbol, qty := range holdings {
		price, err := s.priceService.GetCurrentPrice(symbol)
		if err != nil {
			logrus.Warnf("Failed to get current price for %s: %v", symbol, err)
			continue
		}

		value := price.Mul(qty)
		totalValue = totalValue.Add(value)

		portfolioItems = append(portfolioItems, map[string]interface{}{
			"symbol":        symbol,
			"quantity":      qty,
			"currentPrice":  utils.RoundINR(price),
			"currentValue":  utils.RoundINR(value),
		})
	}

	return map[string]interface{}{
		"userId":      userID,
		"holdings":    portfolioItems,
		"totalValue":  utils.RoundINR(totalValue),
	}, nil
}
