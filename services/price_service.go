package services

import (
	"math/rand"
	"stocky/models"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StockPriceService struct {
	db *gorm.DB
}

func NewStockPriceService(db *gorm.DB) *StockPriceService {
	return &StockPriceService{db: db}
}

// StartPriceUpdater runs hourly to update stock prices
func (s *StockPriceService) StartPriceUpdater() {
	logrus.Info("Starting stock price updater (runs every hour)")

	// Update immediately on startup
	s.updateAllPrices()

	// Then update every hour
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		s.updateAllPrices()
	}
}

// updateAllPrices fetches and updates prices for all stocks
func (s *StockPriceService) updateAllPrices() {
	logrus.Info("Updating stock prices...")

	// Get all unique stock symbols from rewards
	var symbols []string
	if err := s.db.Model(&models.StockReward{}).
		Distinct("stock_symbol").
		Pluck("stock_symbol", &symbols).Error; err != nil {
		logrus.Errorf("Failed to fetch stock symbols: %v", err)
		return
	}

	// Update price for each symbol
	for _, symbol := range symbols {
		price := s.getStockPrice(symbol)

		var stockPrice models.StockPrice
		result := s.db.Where("stock_symbol = ?", symbol).First(&stockPrice)

		if result.Error == gorm.ErrRecordNotFound {
			// Create new price record
			stockPrice = models.StockPrice{
				StockSymbol: symbol,
				Price:       price,
				UpdatedAt:   time.Now(),
			}
			if err := s.db.Create(&stockPrice).Error; err != nil {
				logrus.Errorf("Failed to create price for %s: %v", symbol, err)
			}
		} else {
			// Update existing price record
			stockPrice.Price = price
			stockPrice.UpdatedAt = time.Now()
			if err := s.db.Save(&stockPrice).Error; err != nil {
				logrus.Errorf("Failed to update price for %s: %v", symbol, err)
			}
		}
	}

	logrus.Infof("Updated prices for %d stocks", len(symbols))
}

// getStockPrice generates a hypothetical stock price
func (s *StockPriceService) getStockPrice(symbol string) float64 {
	// Seed based on time for variation
	rand.Seed(time.Now().UnixNano())

	// Base prices for common Indian stocks
	basePrice := map[string]float64{
		"RELIANCE": 2400.0,
		"TCS":      3500.0,
		"INFOSYS":  1500.0,
		"HDFC":     1600.0,
		"ICICI":    900.0,
		"SBI":      600.0,
		"WIPRO":    400.0,
		"BHARTI":   800.0,
		"ITC":      450.0,
		"HCLTECH":  1200.0,
	}

	base, exists := basePrice[symbol]
	if !exists {
		base = 1000.0 // Default base price
	}

	// Add random variation of ±5%
	variation := (rand.Float64() - 0.5) * 0.1 // -5% to +5%
	price := base * (1 + variation)

	return roundToDecimal(price, 2)
}

// GetCurrentPrice returns the current price for a symbol
func (s *StockPriceService) GetCurrentPrice(symbol string) (float64, error) {
	var stockPrice models.StockPrice
	if err := s.db.Where("stock_symbol = ?", symbol).First(&stockPrice).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// If price doesn't exist, generate and save it
			price := s.getStockPrice(symbol)
			stockPrice = models.StockPrice{
				StockSymbol: symbol,
				Price:       price,
				UpdatedAt:   time.Now(),
			}
			if err := s.db.Create(&stockPrice).Error; err != nil {
				return 0, err
			}
			return price, nil
		}
		return 0, err
	}
	return stockPrice.Price, nil
}

func roundToDecimal(value float64, decimals int) float64 {
	multiplier := 1.0
	for i := 0; i < decimals; i++ {
		multiplier *= 10
	}
	return float64(int(value*multiplier+0.5)) / multiplier
}
