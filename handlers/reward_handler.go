package handlers

import (
	"fmt"
	"net/http"
	"stocky/models"
	"stocky/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RewardHandler struct {
	db           *gorm.DB
	priceService *services.StockPriceService
}

func NewRewardHandler(db *gorm.DB) *RewardHandler {
	return &RewardHandler{
		db:           db,
		priceService: services.NewStockPriceService(db),
	}
}

// CreateReward creates a new stock reward
func (h *RewardHandler) CreateReward(c *gin.Context) {
	var req models.RewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default timestamp if not provided
	if req.RewardedAt.IsZero() {
		req.RewardedAt = time.Now()
	}

	// Generate idempotency key if not provided
	if req.IdempotencyKey == "" {
		req.IdempotencyKey = fmt.Sprintf("%s-%s-%d", req.UserID, req.StockSymbol, req.RewardedAt.Unix())
	}

	// Check for duplicate idempotency key
	var existingReward models.StockReward
	if err := h.db.Where("idempotency_key = ?", req.IdempotencyKey).First(&existingReward).Error; err == nil {
		// Duplicate request - return existing reward
		c.JSON(http.StatusConflict, gin.H{
			"error":           "Duplicate reward request",
			"existing_reward": existingReward,
		})
		return
	}

	// Get current stock price
	price, err := h.priceService.GetCurrentPrice(req.StockSymbol)
	if err != nil {
		logrus.Errorf("Failed to get stock price: %v", err)
		price = 0
	}

	// Calculate fees (hypothetical: 0.5% brokerage + 0.1% STT + 18% GST on brokerage)
	stockValue := req.Quantity * price
	brokerage := stockValue * 0.005
	stt := stockValue * 0.001
	gst := brokerage * 0.18
	totalFees := brokerage + stt + gst

	// Begin transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create stock reward
	reward := models.StockReward{
		UserID:         req.UserID,
		StockSymbol:    req.StockSymbol,
		Quantity:       req.Quantity,
		RewardedAt:     req.RewardedAt,
		IdempotencyKey: req.IdempotencyKey,
	}

	if err := tx.Create(&reward).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reward"})
		return
	}

	// Create ledger entries (double-entry bookkeeping)
	ledgerEntries := []models.LedgerEntry{
		{
			RewardID:    reward.ID,
			EntryType:   "STOCK_CREDIT",
			StockSymbol: req.StockSymbol,
			Quantity:    req.Quantity,
			Amount:      stockValue,
			Description: fmt.Sprintf("Stock reward credited to user %s", req.UserID),
		},
		{
			RewardID:    reward.ID,
			EntryType:   "CASH_DEBIT",
			StockSymbol: req.StockSymbol,
			Quantity:    0,
			Amount:      -stockValue,
			Description: fmt.Sprintf("Company cash outflow for stock purchase"),
		},
		{
			RewardID:    reward.ID,
			EntryType:   "FEE_DEBIT",
			StockSymbol: req.StockSymbol,
			Quantity:    0,
			Amount:      -totalFees,
			Description: fmt.Sprintf("Brokerage: %.2f, STT: %.2f, GST: %.2f", brokerage, stt, gst),
		},
	}

	for _, entry := range ledgerEntries {
		if err := tx.Create(&entry).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ledger entry"})
			return
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	logrus.Infof("Created reward for user %s: %f shares of %s", req.UserID, req.Quantity, req.StockSymbol)

	response := models.RewardResponse{
		ID:           reward.ID,
		UserID:       reward.UserID,
		StockSymbol:  reward.StockSymbol,
		Quantity:     reward.Quantity,
		RewardedAt:   reward.RewardedAt,
		CurrentPrice: price,
		CurrentValue: req.Quantity * price,
	}

	c.JSON(http.StatusCreated, response)
}

// GetTodayStocks returns all stock rewards for the user for today
func (h *RewardHandler) GetTodayStocks(c *gin.Context) {
	userID := c.Param("userId")

	// Get today's date range
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var rewards []models.StockReward
	if err := h.db.Where("user_id = ? AND rewarded_at >= ? AND rewarded_at < ?",
		userID, startOfDay, endOfDay).Find(&rewards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rewards"})
		return
	}

	response := models.TodayStocksResponse{
		UserID:  userID,
		Date:    startOfDay.Format("2006-01-02"),
		Rewards: rewards,
	}

	c.JSON(http.StatusOK, response)
}

// GetHistoricalINR returns the INR value of user's stock rewards for all past days
func (h *RewardHandler) GetHistoricalINR(c *gin.Context) {
	userID := c.Param("userId")

	// Get rewards up to yesterday
	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var rewards []models.StockReward
	if err := h.db.Where("user_id = ? AND rewarded_at < ?", userID, startOfToday).
		Order("rewarded_at").Find(&rewards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rewards"})
		return
	}

	// Group by date and calculate INR value
	dailyValues := make(map[string]float64)
	for _, reward := range rewards {
		date := reward.RewardedAt.Format("2006-01-02")
		price, _ := h.priceService.GetCurrentPrice(reward.StockSymbol)
		value := reward.Quantity * price
		dailyValues[date] += value
	}

	// Convert map to slice
	var daily []models.DailyINRValue
	for date, value := range dailyValues {
		daily = append(daily, models.DailyINRValue{
			Date:       date,
			TotalValue: value,
		})
	}

	response := models.HistoricalINRResponse{
		UserID: userID,
		Daily:  daily,
	}

	c.JSON(http.StatusOK, response)
}

// GetStats returns user statistics
func (h *RewardHandler) GetStats(c *gin.Context) {
	userID := c.Param("userId")

	// Get today's rewards
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var todayRewards []models.StockReward
	if err := h.db.Where("user_id = ? AND rewarded_at >= ? AND rewarded_at < ?",
		userID, startOfDay, endOfDay).Find(&todayRewards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch today's rewards"})
		return
	}

	// Group today's rewards by stock symbol
	todayBySymbol := make(map[string]float64)
	for _, reward := range todayRewards {
		todayBySymbol[reward.StockSymbol] += reward.Quantity
	}

	var todayRewardsList []models.StockQuantity
	for symbol, qty := range todayBySymbol {
		todayRewardsList = append(todayRewardsList, models.StockQuantity{
			StockSymbol: symbol,
			Quantity:    qty,
		})
	}

	// Calculate total portfolio value
	var allRewards []models.StockReward
	if err := h.db.Where("user_id = ?", userID).Find(&allRewards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch all rewards"})
		return
	}

	totalValue := 0.0
	for _, reward := range allRewards {
		price, _ := h.priceService.GetCurrentPrice(reward.StockSymbol)
		totalValue += reward.Quantity * price
	}

	response := models.StatsResponse{
		UserID:                userID,
		TodayRewards:          todayRewardsList,
		CurrentPortfolioValue: totalValue,
	}

	c.JSON(http.StatusOK, response)
}

// GetPortfolio shows holdings per stock symbol with current INR value
func (h *RewardHandler) GetPortfolio(c *gin.Context) {
	userID := c.Param("userId")

	var rewards []models.StockReward
	if err := h.db.Where("user_id = ?", userID).Find(&rewards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rewards"})
		return
	}

	// Group by stock symbol
	holdingsBySymbol := make(map[string]float64)
	for _, reward := range rewards {
		holdingsBySymbol[reward.StockSymbol] += reward.Quantity
	}

	// Calculate values
	var holdings []models.HoldingDetail
	totalValue := 0.0

	for symbol, qty := range holdingsBySymbol {
		price, _ := h.priceService.GetCurrentPrice(symbol)
		value := qty * price
		totalValue += value

		holdings = append(holdings, models.HoldingDetail{
			StockSymbol:  symbol,
			TotalShares:  qty,
			CurrentPrice: price,
			CurrentValue: value,
		})
	}

	response := models.PortfolioResponse{
		UserID:     userID,
		Holdings:   holdings,
		TotalValue: totalValue,
	}

	c.JSON(http.StatusOK, response)
}
