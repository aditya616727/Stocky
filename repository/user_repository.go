package repository

import (
	"database/sql"
	"fmt"

	"github.com/aditya616727/stocky/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) getUserByID(userID string) (*models.User, error) {
	query := `
	SELECT id, user_id,name, email, created_at, updated_at
	FROM users
	WHERE user_id = $1	
	`
	var user models.User
	err := r.db.QueryRow(query, userID).Scan(
		&user.ID,
		&user.UserID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
	INSERT INTO users (user_id, name, email) VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(
		query,
		user.UserID, user.Name, user.Email,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) getOrCreateUser(userID string) (*models.User, error) {
	user, err := r.getUserByID(userID)
	if err == nil {
		return user, nil
	}

	newUser := &models.User{
		UserID: userID,
		Name:   "User " + userID,
		Email:  fmt.Sprintf("%s@example.com", userID),
	}
	err = r.CreateUser(newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return newUser, nil
}
