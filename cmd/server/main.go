package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"week5/internal/config"
	"week5/internal/database"
	httpDelivery "week5/internal/delivery/http"
	"week5/internal/repository/postgres"
	"week5/internal/service/email"
	"week5/internal/usecase"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	repo := postgres.NewInventoryRepository(db)
	notifier := email.NewNotifier(cfg)
	invUsecase := usecase.NewInventoryUsecase(repo, notifier)
	handler := httpDelivery.NewInventoryHandler(invUsecase)

	r := gin.Default()
	handler.RegisterRoutes(r)

	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
