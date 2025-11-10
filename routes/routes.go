package routes

import (
	"stocky/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.RouterGroup, db *gorm.DB) {
	handler := handlers.NewRewardHandler(db)

	// Reward endpoints
	router.POST("/reward", handler.CreateReward)
	router.GET("/today-stocks/:userId", handler.GetTodayStocks)
	router.GET("/historical-inr/:userId", handler.GetHistoricalINR)
	router.GET("/stats/:userId", handler.GetStats)
	router.GET("/portfolio/:userId", handler.GetPortfolio)
}
