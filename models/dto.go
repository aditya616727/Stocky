package models

import "time"

type RewardRequest struct {
	userID        string
	StockSymbol   string
	Quantity      float64
	RewardAt      time.Time
	IdempotantKey string
}

type RewardResponse struct {
	Success          bool
	Message          string
	RewardID         string
	UserID           string
	StockSymbol      string
	Quantity         string
	PricePerShare    float64
	TotalValue       float64
	TotalCompanyCost float64
	RewardAt         time.Time
}

type TodayStockResponse struct {
	Success bool
	UserID  string
	Date    string
	Rewards []Reward
	
}
