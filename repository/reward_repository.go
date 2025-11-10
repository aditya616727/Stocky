package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aditya616727/stocky/models"
)

type RewardRepository struct {
	// Implementation goes here
	db *sql.DB
}

func NewRewardRepository(db *sql.DB) *RewardRepository {
	return &RewardRepository{db: db}
}
func (r *RewardRepository) CreateReward(tx *sql.Tx, reward *models.Reward) error {
	// Implementation goes here
	query := `
	INSERT INTO rewards (reward_id, user_id, stock_symbol, quantity, price_per_share, total_cost,brokerage_fee,
	stt_fee,gst_fee,sebi_charges,stap_duty,total_fees,total_company_cost, rewarded_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	RETURNING id, created_at
	`
	err := tx.QueryRow(
		query,
		reward.RewardID, reward.UserID, reward.StockSymbol, reward.Quantity, reward.PricePerShare, reward.TotalCost, reward.BorkerageFee,
		reward.STTPFee, reward.GSTFee, reward.SEBICharges, reward.StampDuty, reward.TotalFees, reward.TotalCompanyCost, reward.RewardedAt,
	).Scan(&reward.ID, &reward.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create reward: %w", err)
	}

	return nil
}
func (r *RewardRepository) GetRewardsByUserAndDate(userID string, date time.Time) ([]models.Reward, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := `
	SELECT 
		id, reward_id, user_id, stock_symbol, quantity, price_per_share, total_cost, brokerage_fee,
		stt_fee, gst_fee, sebi_charges, stap_duty, total_fees, total_company_cost, rewarded_at, created_at
	FROM rewards
	WHERE user_id = $1 AND rewarded_at >= $2 AND rewarded_at < $3
	ORDER BY rewarded_at DESC
	`
	rows, err := r.db.Query(query, userID, startOfDay, endOfDay)
	if err != nil {
		return nil, fmt.Errorf("failed to query rewards: %w", err)
	}
	defer rows.Close()

	var rewards []models.Reward
	for rows.Next() {
		var reward models.Reward
		err := rows.Scan(
			&reward.ID, &reward.RewardID, &reward.UserID, &reward.StockSymbol, &reward.Quantity,
			&reward.PricePerShare, &reward.TotalCost, &reward.BorkerageFee,
			&reward.STTPFee, &reward.GSTFee, &reward.SEBICharges, &reward.StampDuty,
			&reward.TotalFees, &reward.TotalCompanyCost, &reward.RewardedAt, &reward.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reward: %w", err)
		}
		rewards = append(rewards, reward)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return rewards, nil
}

func (r *RewardRepository) GetRewardsByUserBeforeDate(userID string, beforeDate time.Time) ([]models.Reward, error) {
	query := `
	SELECT 
		id, reward_id, user_id, stock_symbol, quantity, price_per_share, total_cost, brokerage_fee,
		stt_fee, gst_fee, sebi_charges, stap_duty, total_fees, total_company_cost, rewarded_at, created_at
	FROM rewards
	WHERE reward_id = $1 AND rewarded_at < $2
	ORDER BY rewarded_at DESC
	`
	rows, err := r.db.Query(query, userID, beforeDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query rewards: %w", err)
	}
	defer rows.Close()
	var rewards []models.Reward
	for rows.Next() {
		var reward models.Reward
		err := rows.Scan(
			&reward.ID, &reward.RewardID, &reward.UserID, &reward.StockSymbol, &reward.Quantity,
			&reward.PricePerShare, &reward.TotalCost, &reward.BorkerageFee,
			&reward.STTPFee, &reward.GSTFee, &reward.SEBICharges, &reward.StampDuty, &reward.TotalFees, &reward.TotalCompanyCost, &reward.RewardedAt, &reward.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reward: %w", err)
		}
		rewards = append(rewards, reward)
	}

	return rewards, nil
}

func (r *RewardRepository) CheckIdempotanceKey(key string) (*models.Idempotancekey, error) {
	query := `
	SELECT id, idempotance_keys, request_payload, response_payload, status,
	FROM idempotance_keys
	WHERE key = $1 and expired_at > NOW()
	`
	var ik models.Idempotancekey
	err := r.db.QueryRow(query, key).Scan(
		&ik.ID, &ik.Idempotancekey, &ik.RequestPayload, &ik.ResponsePayload, &ik.Status,
		&ik.CreatedAt, &ik.ExpiresAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No existing key found
		}
		return nil, fmt.Errorf("failed to check idempotance key: %w", err)
	}
	return &ik, nil
}

func (r *RewardRepository) CreateIdempotanceKey(tx *sql.Tx, key string, requestpayload interface{}, expiresAt time.Time) error {
	payloadJSON, err := json.Marshal(requestpayload)
	if err != nil {
		return fmt.Errorf("failed to marshal request payload: %w", err)
	}
	query := `
	INSERT INTO idempotance_keys (idempotance_keys, request_payload, status, expired_at)
	VALUES ($1, $2, $3, $4)
	`
	_, err = tx.Exec(query, key, string(payloadJSON), models.IdempotancyStatusProcessing, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to create idempotance key: %w", err)
	}
	return nil
}

func (r *RewardRepository) UpdateIdempotanceKey(tx *sql.Tx, key string, responsePayload interface{}, status string) error {
	responseJSON, err := json.Marshal(responsePayload)
	if err != nil {
		return fmt.Errorf("failed to marshal response payload: %w", err)
	}
	query := `
	UPDATE idempotance_keys
	SET response_payload = $1, status = $2,completed_at = NOW()
	WHERE idempotance_keys = $3
	`
	_, err = tx.Exec(query, string(responseJSON), status, key)
	if err != nil {
		return fmt.Errorf("failed to update idempotance key: %w", err)
	}
	return nil
}

func (r *RewardRepository) CleanupExpiredIdempotanceKeys() error {
	query := `
	DELETE FROM idempotance_keys
	WHERE expired_at <= NOW()
	`
	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired idempotance keys: %w", err)
	}
	return nil
}
