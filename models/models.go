package models

import "time"

type User struct {
	ID        int
	UserID    string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Reward struct {
	ID               int
	RewardID         string
	UserID           string
	StockSymbol      string
	Quantity         float64
	PricePerShare    float64
	TotalCost        float64
	BorkerageFee     float64
	STTPFee          float64
	GSTFee           float64
	SEBICharges      float64
	StampDuty        float64
	TotalFees        float64
	TotalCompanyCost float64
	RewardedAt       time.Time
	CreatedAt        time.Time
}

type LedgerEntry struct {
	ID          int
	EntryID     string
	RewardID    string
	AccountType string
	StockSymbol string
	Debit       float64
	Credit      float64
	BalanceType string
	Quantity    *float64
	Description string
	CreatedAt   time.Time
}

const (
	AccountTypeStockAsset     = "STOCK_ASSET"
	AccountTypeCashOutflow    = "CASH_OUTFLOW"
	AccountTypeFeeExpense     = "FEE_EXPENSE"
	AccountTypeStockLiability = "STOCK_LIABLITY"
)

const (
	BalanceTypeStockUnits = "STOCK_UNITS"
	BalanceTypeINR        = "INR"
)

type StockPrice struct {
	ID        int
	Symbol    string
	Price     float64
	Currency  string
	FetchedAt time.Time
	CreatedAt time.Time
}

type UserHolding struct {
	ID            int
	UserId        string
	StockSymbol   string
	TotalQuantity float64
	AverageCost   float64
	TotalCost     float64
	UpdatedAt     time.Time
}

type Idempotancekey struct {
	ID              int
	Idempotancekey  string
	RequestPayload  string
	ResponsePayload string
	Status          string
	CreatedAt       time.Time
	CompletedAt     *time.Time
	ExpiresAt       *time.Time
}

const (
	IdempotancyStatusProcessing = "PROCESSING"
	IdempotancyStatusCompleted  = "COMPLETED"
	IdempotancyStatusFailed     = "FAILED"
)
