package repository

import (
	"database/sql"
	"fmt"

	"github.com/aditya616727/stocky/models"
)

type HoldingRepository struct {
	db *sql.DB
}

func NewHoldingRepository(db *sql.DB) *HoldingRepository {
	return &HoldingRepository{db: db}
}

func (r *HoldingRepository) UpsertHolding(tx *sql.Tx, holding *models.UserHolding) error {
	query := `
    INSERT INTO user_holdings (user_id, stock_symbol, total_quantity, average_cost, total_cost)
    VALUES ($1, $2, $3, $4, $5)
    ON CONFLICT (user_id, stock_symbol)
    DO UPDATE SET
      total_quantity = user_holdings.total_quantity + EXCLUDED.total_quantity,
      total_cost = user_holdings.total_cost + EXCLUDED.total_cost,
      average_cost = (user_holdings.total_cost + EXCLUDED.total_cost) / (user_holdings.total_quantity + EXCLUDED.total_quantity),
      updated_at = CURRENT_TIMESTAMP
    RETURNING id, total_quantity, average_cost, total_cost, updated_at;
    `
	err := tx.QueryRow(
		query,
		holding.UserID,
		holding.StockSymbol,
		holding.TotalQuantity,
		holding.AverageCost,
		holding.TotalCost,
	).Scan(
		&holding.ID,
		&holding.TotalQuantity,
		&holding.AverageCost,
		&holding.TotalCost,
		&holding.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert holding : %w", err)
	}
	return nil
}

func (r *HoldingRepository) GetUserHoldings(UserID string) ([]models.UserHolding, error) {
	query := `
	SELECT id,user_id,stock_symbol, total_quantity,average_cost, total_cost, updated_at
	FROM user_holdings
	WHERE user_id $1 AND total_quantity >1
	ORDER BY stock_symbol
	`

	rows, err := r.db.Query(query, UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to query holding %w", err)
	}
	defer rows.Close()
	var holdings []models.UserHolding
	for rows.Next() {
		var holding models.UserHolding
		err := rows.Scan(
			&holding.ID, &holding.UserID, &holding.StockSymbol, &holding.TotalQuantity, &holding.AverageCost, &holding.TotalCost, &holding.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan holding  : %w", err)

		}
		holdings = append(holdings, holding)

	}
	return holdings, nil

}

func (r *HoldingRepository) GetUserHoldingbyStock(userID, StockSymbol string) (*models.UserHolding, error) {
	query := `
	SELECT id, user_id, stock_symbol, total_quantity,average_cost,total_cost, updated_At
	FROM user_holding
	WHERE user_id =$1 AND stock_symbol =$2
	`
	var holding models.UserHolding
	err := r.db.QueryRow(query, userID, StockSymbol).Scan(
		&holding.ID, &holding.UserID, &holding.StockSymbol,
		&holding.TotalQuantity, &holding.AverageCost, &holding.TotalCost,
		&holding.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get holding %w", err)
	}
	return &holding, nil
}
