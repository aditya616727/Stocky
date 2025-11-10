package main

import (
	"os"

	"github.com/aditya616727/stocky/config"
	"github.com/aditya616727/stocky/database"
	"github.com/aditya616727/stocky/repository"

	"github.com/sirupsen/logrus"
)

func main() {
	initLogger()
	logrus.Info("Starting the Stockey Reward Service . . .")

	//load config
	cfg, err := config.Load()
	if err != nil {
		logrus.Fatalf("failed to load config: %v", err)
	}
	logrus.Infof("loaded configuration for Env: %s", cfg.Server.Enviroment)

	//connect to database
	if err := database.Connect(&cfg.Database); err != nil {
		logrus.Fatal("failed to connect to database ")

	}
	defer database.DB.Close()

	db := database.GetDB()

	//initilize repo
	RewardRepo := repository.NewRewardRepository(db)
	stockPriceRepo := repository.NewStockPriceRepository(db)
	ledgerRepo := repository.NewLedgerRepository(db)
	holdingRepo := repository.NewHoldingRepository(db)
	userRepo := repository.NewUserRepository(db)

	stockPriceService := services.NewStockPriceService(stockPriceRepo)

}

func initLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}
