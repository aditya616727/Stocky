package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/aditya616727/stocky/models"
)

type StockPriceRepository struct {
	db *sql.DB
}

func NewStockPriceRepository(db *sql.DB) *StockPriceRepository {
	return &StockPriceRepository{db: db}
}

func (r *StockPriceRepository) createStockPrice(symbol string, price float64, fetchedAt time.Time) error {
	query := `
	INSERT INTO stock_prices (symbol, price, fetched_at)
	VALUES ($1, $2, $3)
	ON CONFLICT (symbol,fetched_at) DO UPDATE SET price = $2
	`
	_, err := r.db.Exec(query, symbol, price, fetchedAt)
	if err != nil {
		return fmt.Errorf("failed to upsert stock price: %w", err)
	}
	return nil
}
func (r *StockPriceRepository) GetLatestPrice(symbol string) (*models.StockPrice, error) {
	query := `
	SELECT id, symbol, price, currency, fetched_at, created_at
	FROM stock_prices
	WHERE symbol = $1
	ORDER BY fetched_at DESC
	LIMIT 1
	`
	var sp models.StockPrice
	err := r.db.QueryRow(query, symbol).Scan(
		&sp.ID,
		&sp.Symbol,
		&sp.Price,
		&sp.Currency,
		&sp.FetchedAt,
		&sp.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No price found
		}
		return nil, fmt.Errorf("failed to get latest stock price: %w", err)
	}
	return &sp, nil
}
func (r *StockPriceRepository) getPriceAtTime(symbol string, timestamp time.Time) (*models.StockPrice, error) {
	query := `
	SELECT id, symbol, price, currency, fetched_at, created_at
	FROM stock_prices
	WHERE symbol = $1 AND fetched_at <= $2
	ORDER BY fetched_at DESC
	LIMIT 1
	`
	var sp models.StockPrice
	err := r.db.QueryRow(query, symbol, timestamp).Scan(
		&sp.ID,
		&sp.Symbol,
		&sp.Price,
		&sp.Currency,
		&sp.FetchedAt,
		&sp.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No price found
		}
		return nil, fmt.Errorf("failed to get stock price at time: %w", err)
	}
	return &sp, nil
}

func (r *StockPriceRepository) GetAllSymbols() ([]string, error) {
	query := `
	SELECT DISTINCT symbol
	FROM stock_prices
	order by symbol
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query stock symbols: %w", err)
	}
	defer rows.Close()

	var symbols []string
	for rows.Next() {
		var symbol string
		if err := rows.Scan(&symbol); err != nil {
			return nil, fmt.Errorf("failed to scan stock symbol: %w", err)
		}
		symbols = append(symbols, symbol)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return symbols, nil
}

func (r *StockPriceRepository) BulkCreateStockPrices(price map[string]float64, fetchedAt time.Time) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
	INSERT INTO stock_prices (symbol, price, fetched_at)
	VALUES ($1, $2, $3)
	ON CONFLICT (symbol,fetched_at) DO UPDATE SET price = $2
	`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()
	for symbol, price := range price {
		_, err := stmt.Exec(symbol, price, fetchedAt)
		if err != nil {
			return fmt.Errorf("failed to execute statement for symbol %s: %w", symbol, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
