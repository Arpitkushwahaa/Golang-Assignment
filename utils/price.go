package utils

import (
	"fmt"
	"math/rand"
	"stocky-backend/models"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

var rng *rand.Rand

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// PriceGenerator generates mock stock prices
type PriceGenerator struct {
	basePrices map[string]decimal.Decimal
}

// NewPriceGenerator creates a new price generator with base prices
func NewPriceGenerator() *PriceGenerator {
	return &PriceGenerator{
		basePrices: map[string]decimal.Decimal{
			"RELIANCE":    decimal.NewFromFloat(2450.50),
			"TCS":         decimal.NewFromFloat(3680.75),
			"INFY":        decimal.NewFromFloat(1520.30),
			"HDFCBANK":    decimal.NewFromFloat(1650.90),
			"ICICIBANK":   decimal.NewFromFloat(985.25),
			"SBIN":        decimal.NewFromFloat(625.80),
			"BHARTIARTL":  decimal.NewFromFloat(1180.45),
			"ITC":         decimal.NewFromFloat(455.60),
			"KOTAKBANK":   decimal.NewFromFloat(1725.35),
			"LT":          decimal.NewFromFloat(3420.70),
		},
	}
}

// GeneratePrice generates a random price for a stock with variation
func (pg *PriceGenerator) GeneratePrice(symbol string) decimal.Decimal {
	basePrice, exists := pg.basePrices[symbol]
	if !exists {
		// Default base price if symbol not found
		basePrice = decimal.NewFromFloat(1000.0)
		pg.basePrices[symbol] = basePrice
		logrus.Warnf("No base price for %s, using default", symbol)
	}

	// Generate random variation between -5% to +5%
	variation := (rng.Float64() * 10) - 5 // -5 to +5
	variationDecimal := decimal.NewFromFloat(variation / 100)
	
	// Calculate new price: basePrice * (1 + variation)
	newPrice := basePrice.Mul(decimal.NewFromInt(1).Add(variationDecimal))

	// Round to 4 decimal places
	return newPrice.Round(4)
}

// GeneratePricesForAllStocks generates prices for all configured stocks
func (pg *PriceGenerator) GeneratePricesForAllStocks() map[string]decimal.Decimal {
	prices := make(map[string]decimal.Decimal)
	
	for symbol := range pg.basePrices {
		prices[symbol] = pg.GeneratePrice(symbol)
	}
	
	return prices
}

// UpdateBasePrice updates the base price for a symbol (for gradual price movement)
func (pg *PriceGenerator) UpdateBasePrice(symbol string, newPrice decimal.Decimal) {
	pg.basePrices[symbol] = newPrice
}

// CalculateFees calculates brokerage, STT, and GST fees
func CalculateFees(pricePerShare, quantity decimal.Decimal) (brokerage, stt, gst, total decimal.Decimal) {
	totalValue := pricePerShare.Mul(quantity)

	// Brokerage: 0.03% of transaction value or ₹20, whichever is lower
	brokerage = totalValue.Mul(decimal.NewFromFloat(0.0003))
	maxBrokerage := decimal.NewFromInt(20)
	if brokerage.GreaterThan(maxBrokerage) {
		brokerage = maxBrokerage
	}

	// STT (Securities Transaction Tax): 0.1% on buy side
	stt = totalValue.Mul(decimal.NewFromFloat(0.001))

	// GST: 18% on brokerage
	gst = brokerage.Mul(decimal.NewFromFloat(0.18))

	// Total fees
	total = brokerage.Add(stt).Add(gst)

	// Round all to 4 decimal places
	brokerage = brokerage.Round(4)
	stt = stt.Round(4)
	gst = gst.Round(4)
	total = total.Round(4)

	return
}

// ValidateQuantity ensures quantity is positive and within limits
func ValidateQuantity(quantity decimal.Decimal) error {
	if quantity.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("quantity must be positive")
	}

	// Maximum 10000 shares per transaction
	maxQuantity := decimal.NewFromInt(10000)
	if quantity.GreaterThan(maxQuantity) {
		return fmt.Errorf("quantity cannot exceed 10000 shares")
	}

	return nil
}

// ValidateStockSymbol validates stock symbol format
func ValidateStockSymbol(symbol string) error {
	if len(symbol) == 0 || len(symbol) > 20 {
		return fmt.Errorf("invalid stock symbol length")
	}
	return nil
}

// RoundINR rounds INR amount to 4 decimal places
func RoundINR(amount decimal.Decimal) decimal.Decimal {
	return amount.Round(4)
}

// RoundQuantity rounds quantity to 6 decimal places
func RoundQuantity(quantity decimal.Decimal) decimal.Decimal {
	return quantity.Round(6)
}
