package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	StockPrice StockPriceConfig
	fees       FeesConfig
}

type ServerConfig struct {
	Port       string
	Host       string
	Enviroment string
}

type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	Name         string
	SSLMode      string
	MaxOpenConns int
	MaxIdConns   int
}

type StockPriceConfig struct {
	UpdateInterVal       string
	EnablePriceSchedular bool
}

type FeesConfig struct {
	BrokerageFeePercent   float64
	STTFeePercent         float64
	GSTOnBrokeragePercent float64
	SEBIChargesPercent    float64
	StampDutyPercent      float64
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("no .env file found")
	}
	MaxOpenConns, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
	MaxIdConns, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "5"))

	BrokerageFee, _ := strconv.ParseFloat(getEnv("BROKERAGE_FEE_PERCENTAGE", "0.03"), 64)
	sttFee, _ := strconv.ParseFloat(getEnv("STT_FEE_PERCENT", "0.1"), 64)
	gstOnBrokerage, _ := strconv.ParseFloat(getEnv("GST_ON_BROKERAGE_PERCENT", "18"), 64)
	SEBICharges, _ := strconv.ParseFloat(getEnv("SEBI_CHARGES_PERCENT", "0.0001"), 64)
	stampDuty, _ := strconv.ParseFloat(getEnv("STAMP_DUTY_PERCENT", "0.015"), 64)

	enableSchedular, _ := strconv.ParseBool(getEnv("ENABLE_PRICE_SCHEDULAR", "true"))

	Config := &Config{
		Server: ServerConfig{
			Port:       getEnv("SERVER_PORT", "8080"),
			Host:       getEnv("SERVER_HOST", "localhost"),
			Enviroment: getEnv("ENVIROMENT", "developemnt"),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "5432"),
			User:         getEnv("DB_USER", "postgre"),
			Password:     getEnv("DB_PASSWORD", "postgre"),
			Name:         getEnv("DB_NAME", "assignment"),
			SSLMode:      getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns: MaxOpenConns,
			MaxIdConns:   MaxIdConns,
		},
		StockPrice: StockPriceConfig{
			UpdateInterVal:       getEnv("STOCK_PRICE_UPDATE_INTERVAL", "1h"),
			EnablePriceSchedular: enableSchedular,
		},
		fees: FeesConfig{
			BrokerageFeePercent:   BrokerageFee,
			STTFeePercent:         sttFee,
			GSTOnBrokeragePercent: gstOnBrokerage,
			SEBIChargesPercent:    SEBICharges,
			StampDutyPercent:      stampDuty,
		},
	}
	return Config, nil

}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}
