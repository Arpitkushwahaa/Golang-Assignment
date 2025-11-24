package services

import (
	"fmt"
	"stocky-backend/db"
	"stocky-backend/models"
	"stocky-backend/utils"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// PriceService handles stock price operations
type PriceService struct {
	generator *utils.PriceGenerator
}

// NewPriceService creates a new price service
func NewPriceService() *PriceService {
	return &PriceService{
		generator: utils.NewPriceGenerator(),
	}
}

// GetCurrentPrice retrieves the current price for a stock
func (s *PriceService) GetCurrentPrice(symbol string) (decimal.Decimal, error) {
	var priceHistory models.PriceHistory

	// Try to get the latest price from database
	err := db.DB.Where("stock_symbol = ?", symbol).
		Order("timestamp DESC").
		First(&priceHistory).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Generate and save a new price
			price := s.generator.GeneratePrice(symbol)
			if err := s.SavePrice(symbol, price); err != nil {
				logrus.Warnf("Failed to save generated price: %v", err)
			}
			return price, nil
		}
		return decimal.Zero, fmt.Errorf("failed to fetch price: %w", err)
	}

	// Check if price is recent (within last 2 hours)
	if time.Since(priceHistory.Timestamp) > 2*time.Hour {
		logrus.Warnf("Price for %s is stale, generating new price", symbol)
		price := s.generator.GeneratePrice(symbol)
		if err := s.SavePrice(symbol, price); err != nil {
			logrus.Warnf("Failed to save generated price: %v", err)
		}
		return price, nil
	}

	return priceHistory.PriceINR, nil
}

// GetPriceAtTime retrieves the price for a stock at a specific time
func (s *PriceService) GetPriceAtTime(symbol string, timestamp time.Time) (decimal.Decimal, error) {
	var priceHistory models.PriceHistory

	// Get the closest price before or at the given timestamp
	err := db.DB.Where("stock_symbol = ? AND timestamp <= ?", symbol, timestamp).
		Order("timestamp DESC").
		First(&priceHistory).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// No historical price found, use current price as fallback
			logrus.Warnf("No price history found for %s at %v, using current price", symbol, timestamp)
			return s.GetCurrentPrice(symbol)
		}
		return decimal.Zero, fmt.Errorf("failed to fetch historical price: %w", err)
	}

	return priceHistory.PriceINR, nil
}

// SavePrice saves a new price to the database
func (s *PriceService) SavePrice(symbol string, price decimal.Decimal) error {
	priceHistory := models.PriceHistory{
		StockSymbol: symbol,
		PriceINR:    utils.RoundINR(price),
		Timestamp:   utils.NowUTC(),
	}

	if err := db.DB.Create(&priceHistory).Error; err != nil {
		return fmt.Errorf("failed to save price: %w", err)
	}

	// Update generator base price for gradual movement
	s.generator.UpdateBasePrice(symbol, price)

	logrus.WithFields(logrus.Fields{
		"symbol": symbol,
		"price":  price,
	}).Debug("Price saved successfully")

	return nil
}

// UpdateAllPrices generates and saves new prices for all stocks
func (s *PriceService) UpdateAllPrices() error {
	logrus.Info("Updating prices for all stocks...")

	prices := s.generator.GeneratePricesForAllStocks()
	
	for symbol, price := range prices {
		if err := s.SavePrice(symbol, price); err != nil {
			logrus.Errorf("Failed to save price for %s: %v", symbol, err)
			continue
		}
	}

	logrus.Info("All prices updated successfully")
	return nil
}

// GetPricesForDate retrieves all prices for a specific date
func (s *PriceService) GetPricesForDate(date time.Time) (map[string]decimal.Decimal, error) {
	startOfDay := utils.StartOfDayUTC(date)
	endOfDay := utils.EndOfDayUTC(date)

	var prices []models.PriceHistory

	// Get distinct latest prices for each symbol on that day
	err := db.DB.Raw(`
		SELECT DISTINCT ON (stock_symbol) *
		FROM price_history
		WHERE timestamp >= ? AND timestamp <= ?
		ORDER BY stock_symbol, timestamp DESC
	`, startOfDay, endOfDay).Scan(&prices).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch prices for date: %w", err)
	}

	priceMap := make(map[string]decimal.Decimal)
	for _, p := range prices {
		priceMap[p.StockSymbol] = p.PriceINR
	}

	return priceMap, nil
}
