package services

import (
	"fmt"
	"stocky-backend/db"
	"stocky-backend/models"
	"stocky-backend/utils"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

// LedgerService handles ledger entry operations
type LedgerService struct{}

// NewLedgerService creates a new ledger service
func NewLedgerService() *LedgerService {
	return &LedgerService{}
}

// CreateLedgerEntries creates double-entry accounting entries for a reward
func (s *LedgerService) CreateLedgerEntries(rewardEvent *models.RewardEvent, pricePerShare decimal.Decimal) error {
	quantity := rewardEvent.Quantity
	symbol := rewardEvent.StockSymbol
	timestamp := rewardEvent.Timestamp

	// Calculate total stock value
	totalValue := pricePerShare.Mul(quantity)
	totalValue = utils.RoundINR(totalValue)

	// Calculate fees
	brokerage, stt, gst, totalFees := utils.CalculateFees(pricePerShare, quantity)

	logrus.WithFields(logrus.Fields{
		"rewardEventId": rewardEvent.ID,
		"symbol":        symbol,
		"quantity":      quantity,
		"pricePerShare": pricePerShare,
		"totalValue":    totalValue,
		"brokerage":     brokerage,
		"stt":           stt,
		"gst":           gst,
		"totalFees":     totalFees,
	}).Info("Creating ledger entries")

	entries := []models.LedgerEntry{
		{
			// STOCK entry: +X shares credited to user
			RewardEventID: rewardEvent.ID,
			EntryType:     models.EntryTypeStock,
			StockSymbol:   &symbol,
			Quantity:      quantity,
			AmountINR:     totalValue,
			Timestamp:     timestamp,
		},
		{
			// CASH entry: Company pays for stocks
			RewardEventID: rewardEvent.ID,
			EntryType:     models.EntryTypeCash,
			StockSymbol:   &symbol,
			Quantity:      decimal.Zero,
			AmountINR:     totalValue.Neg(), // Negative for company outflow
			Timestamp:     timestamp,
		},
		{
			// FEE entry: Company pays brokerage
			RewardEventID: rewardEvent.ID,
			EntryType:     models.EntryTypeFee,
			StockSymbol:   &symbol,
			Quantity:      decimal.Zero,
			AmountINR:     totalFees.Neg(), // Negative for company outflow
			Timestamp:     timestamp,
		},
	}

	// Create all entries in a transaction
	if err := db.DB.Create(&entries).Error; err != nil {
		return fmt.Errorf("failed to create ledger entries: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"rewardEventId": rewardEvent.ID,
		"entriesCount":  len(entries),
	}).Info("Ledger entries created successfully")

	return nil
}

// GetUserStockHoldings retrieves total stock holdings for a user
func (s *LedgerService) GetUserStockHoldings(userID int) (map[string]decimal.Decimal, error) {
	type HoldingResult struct {
		StockSymbol string
		TotalQty    decimal.Decimal
	}

	var holdings []HoldingResult

	// Sum up all STOCK entries for the user (group by symbol)
	err := db.DB.Raw(`
		SELECT le.stock_symbol, SUM(le.quantity) as total_qty
		FROM ledger_entries le
		JOIN reward_events re ON le.reward_event_id = re.id
		WHERE re.user_id = ? 
		  AND le.entry_type = 'STOCK'
		  AND le.stock_symbol IS NOT NULL
		  AND re.deleted_at IS NULL
		GROUP BY le.stock_symbol
		HAVING SUM(le.quantity) > 0
	`, userID).Scan(&holdings).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch holdings: %w", err)
	}

	holdingsMap := make(map[string]decimal.Decimal)
	for _, h := range holdings {
		holdingsMap[h.StockSymbol] = h.TotalQty
	}

	return holdingsMap, nil
}

// GetUserStockHoldingsUpToDate retrieves holdings up to a specific date
func (s *LedgerService) GetUserStockHoldingsUpToDate(userID int, endDate time.Time) (map[string]decimal.Decimal, error) {
	type HoldingResult struct {
		StockSymbol string
		TotalQty    decimal.Decimal
	}

	var holdings []HoldingResult

	// Sum up all STOCK entries for the user up to the given date
	err := db.DB.Raw(`
		SELECT le.stock_symbol, SUM(le.quantity) as total_qty
		FROM ledger_entries le
		JOIN reward_events re ON le.reward_event_id = re.id
		WHERE re.user_id = ? 
		  AND le.entry_type = 'STOCK'
		  AND le.stock_symbol IS NOT NULL
		  AND re.deleted_at IS NULL
		  AND re.timestamp <= ?
		GROUP BY le.stock_symbol
		HAVING SUM(le.quantity) > 0
	`, userID, endDate).Scan(&holdings).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch holdings: %w", err)
	}

	holdingsMap := make(map[string]decimal.Decimal)
	for _, h := range holdings {
		holdingsMap[h.StockSymbol] = h.TotalQty
	}

	return holdingsMap, nil
}

// CreateReversalEntries creates reversal ledger entries for reward cancellation
func (s *LedgerService) CreateReversalEntries(rewardEvent *models.RewardEvent) error {
	// Get original ledger entries
	var originalEntries []models.LedgerEntry
	err := db.DB.Where("reward_event_id = ?", rewardEvent.ID).Find(&originalEntries).Error
	if err != nil {
		return fmt.Errorf("failed to fetch original entries: %w", err)
	}

	// Create reversal entries with negative amounts
	var reversalEntries []models.LedgerEntry
	for _, entry := range originalEntries {
		reversal := models.LedgerEntry{
			RewardEventID: entry.RewardEventID,
			EntryType:     entry.EntryType,
			StockSymbol:   entry.StockSymbol,
			Quantity:      entry.Quantity.Neg(),
			AmountINR:     entry.AmountINR.Neg(),
			Timestamp:     utils.NowUTC(),
		}
		reversalEntries = append(reversalEntries, reversal)
	}

	if err := db.DB.Create(&reversalEntries).Error; err != nil {
		return fmt.Errorf("failed to create reversal entries: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"rewardEventId":   rewardEvent.ID,
		"reversalEntries": len(reversalEntries),
	}).Info("Reversal entries created successfully")

	return nil
}
