package main

import (
	"stocky/config"
	"stocky/database"
	"stocky/routes"
	"stocky/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Setup logger
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		logrus.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize stock price service
	priceService := services.NewStockPriceService(db)
	go priceService.StartPriceUpdater() // Start hourly price updates

	// Setup router
	router := gin.Default()

	// API routes
	api := router.Group("/api/v1")
	routes.SetupRoutes(api, db)

	// Start server
	logrus.Info("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
