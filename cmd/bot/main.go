package main

import (
	"context"
	"log"
	"os"

	"github.com/crocxdued/currency-telegram-bot/internal/app"
	"github.com/crocxdued/currency-telegram-bot/internal/config"
	"github.com/crocxdued/currency-telegram-bot/pkg/logger"
	"github.com/pressly/goose/v3"
)

func main() {
	log.Println("🎯 BOT MAIN STARTED")

	// Загружаем конфигурацию
	log.Println("📋 Loading configuration...")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}
	log.Println("✅ Configuration loaded")

	// Инициализируем логгер
	log.Println("📝 Initializing logger...")
	if err := logger.InitGlobal(cfg.LogLevel); err != nil {
		log.Fatalf("❌ Failed to initialize logger: %v", err)
	}
	logger.S.Info("✅ Logger initialized")

	// Обработка команды миграций
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		logger.S.Info("🔄 Running migrations...")
		if err := runMigrations(cfg); err != nil {
			logger.S.Fatalf("❌ Migration failed: %v", err)
		}
		logger.S.Info("✅ Migrations completed")
		return
	}

	// Создаем и запускаем приложение
	logger.S.Info("🚀 Creating application...")
	application := app.New(cfg)

	logger.S.Info("🎯 Starting application...")
	if err := application.Run(); err != nil {
		logger.S.Errorf("❌ Application failed: %v", err)
		os.Exit(1)
	}
}

func runMigrations(cfg *config.Config) error {
	logger.S.Infof("🔗 Connecting to database: %s", cfg.DBURL[:30]+"...")
	db, err := goose.OpenDBWithDriver("postgres", cfg.DBURL)
	if err != nil {
		return err
	}
	defer db.Close()

	command := "up"
	if len(os.Args) > 2 {
		command = os.Args[2]
	}

	logger.S.Infof("🔄 Running migration command: %s", command)
	return goose.RunContext(context.Background(), command, db, "migrations")
}
