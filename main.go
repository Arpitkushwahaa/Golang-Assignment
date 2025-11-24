package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"stocky-backend/db"
	"stocky-backend/routes"
	"stocky-backend/services"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logger
	initLogger()

	logrus.Info("Starting Stocky Backend...")

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, using environment variables")
	}

	// Initialize database
	if err := db.Initialize(); err != nil {
		logrus.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Start price update scheduler
	priceService := services.NewPriceService()
	startPriceUpdateScheduler(priceService)

	// Setup router
	router := routes.SetupRouter()

	// Get server port
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logrus.Infof("Server is running on http://localhost:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")

	// Graceful shutdown with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatalf("Server forced to shutdown: %v", err)
	}

	logrus.Info("Server exited successfully")
}

// initLogger initializes the logrus logger
func initLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	logrus.SetOutput(os.Stdout)

	// Set log level based on environment
	if os.Getenv("GIN_MODE") == "release" {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

// startPriceUpdateScheduler starts a background scheduler to update prices hourly
func startPriceUpdateScheduler(priceService *services.PriceService) {
	logrus.Info("Starting price update scheduler (hourly)")

	// Update prices immediately on startup
	go func() {
		if err := priceService.UpdateAllPrices(); err != nil {
			logrus.Errorf("Failed to update prices on startup: %v", err)
		}
	}()

	// Schedule hourly updates
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			logrus.Info("Scheduled price update triggered")
			if err := priceService.UpdateAllPrices(); err != nil {
				logrus.Errorf("Failed to update prices: %v", err)
			}
		}
	}()
}
