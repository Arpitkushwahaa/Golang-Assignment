package db

import (
	"fmt"
	"os"
	"stocky-backend/models"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Initialize connects to PostgreSQL and runs migrations
func Initialize() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	logrus.Info("Connecting to database...")

	// Configure GORM logger
	gormLogger := logger.Default.LogMode(logger.Info)
	if os.Getenv("GIN_MODE") == "release" {
		gormLogger = logger.Default.LogMode(logger.Error)
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	logrus.Info("Database connected successfully")

	// Run migrations
	if err := runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// runMigrations creates all required tables
func runMigrations() error {
	logrus.Info("Running database migrations...")

	// Auto-migrate all models
	err := DB.AutoMigrate(
		&models.RewardEvent{},
		&models.LedgerEntry{},
		&models.PriceHistory{},
		&models.StockConfig{},
	)
	if err != nil {
		return err
	}

	// Create composite indexes for better query performance
	if err := createIndexes(); err != nil {
		return err
	}

	// Initialize stock configurations
	if err := initializeStockConfigs(); err != nil {
		return err
	}

	logrus.Info("Database migrations completed successfully")
	return nil
}

// createIndexes creates additional composite indexes
func createIndexes() error {
	logrus.Info("Creating composite indexes...")

	// Composite index for reward events by user and date
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_reward_user_timestamp ON reward_events(user_id, timestamp DESC)")

	// Composite index for price history lookup
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_price_symbol_timestamp ON price_history(stock_symbol, timestamp DESC)")

	// Index for ledger entries by type
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_ledger_entry_type ON ledger_entries(entry_type)")

	// Unique constraint for deduplication
	DB.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_reward_dedup 
		ON reward_events(user_id, stock_symbol, quantity, timestamp) 
		WHERE deleted_at IS NULL
	`)

	return nil
}

// initializeStockConfigs seeds initial stock configurations
func initializeStockConfigs() error {
	logrus.Info("Initializing stock configurations...")

	stocks := []string{"RELIANCE", "TCS", "INFY", "HDFCBANK", "ICICIBANK", "SBIN", "BHARTIARTL", "ITC", "KOTAKBANK", "LT"}

	for _, symbol := range stocks {
		var count int64
		DB.Model(&models.StockConfig{}).Where("stock_symbol = ?", symbol).Count(&count)

		if count == 0 {
			config := models.StockConfig{
				StockSymbol: symbol,
				Multiplier:  mustParseDecimal("1.0"),
				IsActive:    true,
				Notes:       "Initial configuration",
			}
			if err := DB.Create(&config).Error; err != nil {
				logrus.Warnf("Failed to create stock config for %s: %v", symbol, err)
			}
		}
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// Helper function to parse decimal (panics on error - only for init)
func mustParseDecimal(s string) decimal.Decimal {
	d, err := decimal.NewFromString(s)
	if err != nil {
		panic(err)
	}
	return d
}
