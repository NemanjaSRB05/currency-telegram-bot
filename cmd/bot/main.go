package main

import (
	"log"
	"os"

	"github.com/crocxdued/currency-telegram-bot/internal/app"
	"github.com/crocxdued/currency-telegram-bot/internal/config"
	"github.com/crocxdued/currency-telegram-bot/pkg/logger"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализируем логгер
	if err := logger.InitGlobal(cfg.LogLevel); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Создаем и запускаем приложение
	application := app.New(cfg)
	if err := application.Run(); err != nil {
		logger.S.Errorf("Application failed: %v", err)
		os.Exit(1)
	}
}
