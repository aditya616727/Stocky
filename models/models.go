package models

import (
	"time"
)

// StockReward represents a reward event
type StockReward struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	UserID         string    `json:"user_id" gorm:"type:varchar(100);not null;index"`
	StockSymbol    string    `json:"stock_symbol" gorm:"type:varchar(20);not null;index"`
	Quantity       float64   `json:"quantity" gorm:"type:numeric(18,6);not null"`
	RewardedAt     time.Time `json:"rewarded_at" gorm:"not null;index"`
	IdempotencyKey string    `json:"idempotency_key" gorm:"type:varchar(100);uniqueIndex"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// LedgerEntry represents double-entry bookkeeping
type LedgerEntry struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	RewardID    int64     `json:"reward_id" gorm:"not null;index"`
	EntryType   string    `json:"entry_type" gorm:"type:varchar(50);not null"` // STOCK_CREDIT, CASH_DEBIT, FEE_DEBIT
	StockSymbol string    `json:"stock_symbol" gorm:"type:varchar(20)"`
	Quantity    float64   `json:"quantity" gorm:"type:numeric(18,6)"`
	Amount      float64   `json:"amount" gorm:"type:numeric(18,4)"` // INR amount
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
}

// StockPrice represents current stock prices
type StockPrice struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	StockSymbol string    `json:"stock_symbol" gorm:"type:varchar(20);uniqueIndex;not null"`
	Price       float64   `json:"price" gorm:"type:numeric(18,4);not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"index"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserPortfolio represents aggregated user holdings
type UserPortfolio struct {
	UserID       string  `json:"user_id"`
	StockSymbol  string  `json:"stock_symbol"`
	TotalShares  float64 `json:"total_shares"`
	CurrentPrice float64 `json:"current_price"`
	CurrentValue float64 `json:"current_value"`
}

// RewardRequest API request for creating reward
type RewardRequest struct {
	UserID         string    `json:"user_id" binding:"required" example:"user123"`
	StockSymbol    string    `json:"stock_symbol" binding:"required" example:"RELIANCE"`
	Quantity       float64   `json:"quantity" binding:"required,gt=0" example:"10.5"`
	RewardedAt     time.Time `json:"rewarded_at" example:"2025-11-10T10:00:00Z"`
	IdempotencyKey string    `json:"idempotency_key" example:"reward-123-456"`
}

// RewardResponse API response for reward creation
type RewardResponse struct {
	ID           int64     `json:"id" example:"1"`
	UserID       string    `json:"user_id" example:"user123"`
	StockSymbol  string    `json:"stock_symbol" example:"RELIANCE"`
	Quantity     float64   `json:"quantity" example:"10.5"`
	RewardedAt   time.Time `json:"rewarded_at" example:"2025-11-10T10:00:00Z"`
	CurrentPrice float64   `json:"current_price" example:"2450.75"`
	CurrentValue float64   `json:"current_value" example:"25732.88"`
}

// TodayStocksResponse API response for today's stocks
type TodayStocksResponse struct {
	UserID  string        `json:"user_id" example:"user123"`
	Date    string        `json:"date" example:"2025-11-10"`
	Rewards []StockReward `json:"rewards"`
}

// HistoricalINRResponse API response for historical INR values
type HistoricalINRResponse struct {
	UserID string          `json:"user_id" example:"user123"`
	Daily  []DailyINRValue `json:"daily"`
}

// DailyINRValue represents INR value for a specific day
type DailyINRValue struct {
	Date       string  `json:"date" example:"2025-11-09"`
	TotalValue float64 `json:"total_value" example:"125000.50"`
}

// StatsResponse API response for user stats
type StatsResponse struct {
	UserID                string          `json:"user_id" example:"user123"`
	TodayRewards          []StockQuantity `json:"today_rewards"`
	CurrentPortfolioValue float64         `json:"current_portfolio_value" example:"250000.75"`
}

// StockQuantity represents quantity by stock symbol
type StockQuantity struct {
	StockSymbol string  `json:"stock_symbol" example:"RELIANCE"`
	Quantity    float64 `json:"quantity" example:"15.5"`
}

// PortfolioResponse API response for portfolio
type PortfolioResponse struct {
	UserID     string          `json:"user_id" example:"user123"`
	Holdings   []HoldingDetail `json:"holdings"`
	TotalValue float64         `json:"total_value" example:"250000.75"`
}

// HoldingDetail represents individual stock holding
type HoldingDetail struct {
	StockSymbol  string  `json:"stock_symbol" example:"RELIANCE"`
	TotalShares  float64 `json:"total_shares" example:"25.5"`
	CurrentPrice float64 `json:"current_price" example:"2450.75"`
	CurrentValue float64 `json:"current_value" example:"62489.13"`
}
