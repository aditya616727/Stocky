package database

import (
	"database/sql"
	"fmt"

	"github.com/aditya616727/stocky/config"
	"github.com/sirupsen/logrus"
)

var DB *sql.DB

func Connect(cfg *config.DatabaseConfig) error {
	var err error

	dsn := cfg.GetDSN()
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database conn. %w", err)
	}
	//set connection pool setting
	DB.SetMaxOpenConns(cfg.MaxOpenConns)
	DB.SetMaxIdleConns(cfg.MaxIdConns)
	//verify conn..
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping db")
	}
	logrus.Info("connected")
	return nil
}

func getDSN() *sql.DB {
	return DB
}
