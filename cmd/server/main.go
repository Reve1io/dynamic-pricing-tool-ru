package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"dynamic-pricing-tool-ru/internal/api"
	"dynamic-pricing-tool-ru/internal/config"
	"dynamic-pricing-tool-ru/internal/logger"
	"dynamic-pricing-tool-ru/internal/processor"
	"dynamic-pricing-tool-ru/internal/server"
)

func main() {
	if err := logger.Init("logs"); err != nil {
		panic(err)
	}
	defer logger.L.Sync()

	if err := godotenv.Load(); err != nil {
		logger.L.Info("No .env file found, using environment variables")
	}

	cfg := config.LoadConfig()

	if cfg.GetchipsToken == "" {
		logger.L.Fatal("GETCHIPS_TOKEN is required. Set it in .env file or environment variable")
	}

	if cfg.EfindToken == "" {
		logger.L.Fatal("EFIND_TOKEN is required. Set it in .env file or environment variable")
	}

	getchipsClient := api.NewGetchipsClient(cfg.GetchipsURL, cfg.GetchipsToken)
	efindClient := api.NewEfindClient(cfg.EfindURL, cfg.EfindToken)

	proc := processor.NewProcessorWithClients(getchipsClient, efindClient, cfg.ChunkSize)

	handler := server.NewHandler(proc)

	router := gin.Default()

	router.Use(logger.RequestID())
	router.Use(logger.GinLogger())
	router.Use(gin.Recovery())

	router.POST("/api/v1/ru/process", handler.HandleProcess)
	router.GET("/health", handler.HealthCheck)

	logger.L.Info("Server starting",
		zap.String("port", cfg.ServerPort),
	)

	if err := router.Run(":" + cfg.ServerPort); err != nil {
		logger.L.Fatal("Failed to start server:",
			zap.Error(err))
	}
}
