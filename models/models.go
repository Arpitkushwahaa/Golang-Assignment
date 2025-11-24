package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// RewardEvent represents a stock reward given to a user
type RewardEvent struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	UserID      int             `gorm:"not null;index:idx_user_rewards" json:"userId"`
	StockSymbol string          `gorm:"not null;size:20;index:idx_stock_symbol" json:"symbol"`
	Quantity    decimal.Decimal `gorm:"type:numeric(18,6);not null" json:"quantity"`
	Timestamp   time.Time       `gorm:"not null;index:idx_timestamp" json:"timestamp"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`

	// Relationships
	LedgerEntries []LedgerEntry `gorm:"foreignKey:RewardEventID" json:"-"`
}

// TableName specifies the table name for RewardEvent
func (RewardEvent) TableName() string {
	return "reward_events"
}

// EntryType represents the type of ledger entry
type EntryType string

const (
	EntryTypeStock EntryType = "STOCK"
	EntryTypeCash  EntryType = "CASH"
	EntryTypeFee   EntryType = "FEE"
)

// LedgerEntry represents double-entry accounting for rewards
type LedgerEntry struct {
	ID             uint            `gorm:"primaryKey" json:"id"`
	RewardEventID  uint            `gorm:"not null;index:idx_reward_event" json:"rewardEventId"`
	EntryType      EntryType       `gorm:"type:varchar(10);not null" json:"entryType"`
	StockSymbol    *string         `gorm:"size:20" json:"stockSymbol,omitempty"`
	Quantity       decimal.Decimal `gorm:"type:numeric(18,6);not null;default:0" json:"quantity"`
	AmountINR      decimal.Decimal `gorm:"type:numeric(18,4);not null;default:0" json:"amountInr"`
	Timestamp      time.Time       `gorm:"not null" json:"timestamp"`
	CreatedAt      time.Time       `json:"createdAt"`

	// Relationships
	RewardEvent RewardEvent `gorm:"foreignKey:RewardEventID" json:"-"`
}

// TableName specifies the table name for LedgerEntry
func (LedgerEntry) TableName() string {
	return "ledger_entries"
}

// PriceHistory stores historical stock prices
type PriceHistory struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	StockSymbol string          `gorm:"not null;size:20;index:idx_symbol_time" json:"stockSymbol"`
	PriceINR    decimal.Decimal `gorm:"type:numeric(18,4);not null" json:"priceInr"`
	Timestamp   time.Time       `gorm:"not null;index:idx_symbol_time" json:"timestamp"`
	CreatedAt   time.Time       `json:"createdAt"`
}

// TableName specifies the table name for PriceHistory
func (PriceHistory) TableName() string {
	return "price_history"
}

// StockConfig stores configuration for stocks (splits, multipliers, etc.)
type StockConfig struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	StockSymbol string          `gorm:"not null;size:20;uniqueIndex" json:"stockSymbol"`
	Multiplier  decimal.Decimal `gorm:"type:numeric(18,6);not null;default:1" json:"multiplier"`
	IsActive    bool            `gorm:"not null;default:true" json:"isActive"`
	Notes       string          `gorm:"type:text" json:"notes,omitempty"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	CreatedAt   time.Time       `json:"createdAt"`
}

// TableName specifies the table name for StockConfig
func (StockConfig) TableName() string {
	return "stock_config"
}
