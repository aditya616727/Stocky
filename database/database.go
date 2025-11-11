package database

import (
	"fmt"
	"stocky/config"
	"stocky/models"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	var dsn string
	if cfg.DBPassword == "" {
		dsn = fmt.Sprintf(
			"host=%s user=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
			cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort,
		)
	} else {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
		)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	logrus.Info("Connected to database successfully")
	return db, nil
}

func RunMigrations(db *gorm.DB) error {
	logrus.Info("Running database migrations...")

	err := db.AutoMigrate(
		&models.StockReward{},
		&models.LedgerEntry{},
		&models.StockPrice{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logrus.Info("Database migrations completed successfully")
	return nil
}
