package repository

import (
	"database/sql"
	"fmt"

	"github.com/aditya616727/stocky/models"
)

type LedgerRepository struct {
	db *sql.DB
}

func NewLedgerRepository(db *sql.DB) *LedgerRepository {
	return &LedgerRepository{db: db}
}

func (r *LedgerRepository) CreateledgerEntry(tx *sql.Tx, entry *models.LedgerEntry) error {
	query := `
    INSERT INTO ledger_entries (
        entry_id, 
        reward_id, 
        account_type, 
        stock_symbol,
        debit,
        credit,
        balance_type, 
        quantity,
        description
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    RETURNING id, created_at
`
	err := tx.QueryRow(
		query,
		entry.EntryID, entry.RewardID, entry.AccountType, entry.StockSymbol,
		entry.Debit, entry.CreatedAt, entry.BalanceType, entry.Quantity, entry.Description,
	).Scan(&entry.ID, entry.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create ledger entry: %w", err)
	}
	return nil
}

func (r *LedgerRepository) GetLedgerEntriesByRewardID(rewardID string) ([]models.LedgerEntry, error) {
	query := `
	SELECT 
		id, entry_id, reward_id, account_type, stock_symbol, debit, credit, balance_type, quantity, description, created_at
	FROM ledger_entries
	WHERE reward_id = $1
	ORDER BY created_at ASC
	`
	rows, err := r.db.Query(query, rewardID)
	if err != nil {
		return nil, fmt.Errorf("failed to query ledger entries: %w", err)
	}
	defer rows.Close()

	var entries []models.LedgerEntry
	for rows.Next() {
		var entry models.LedgerEntry
		err := rows.Scan(&entry.ID, &entry.EntryID, &entry.RewardID, &entry.AccountType, &entry.StockSymbol, &entry.Debit, &entry.Credit, &entry.BalanceType, &entry.Quantity, &entry.Description, &entry.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ledger entry: %w", err)
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
