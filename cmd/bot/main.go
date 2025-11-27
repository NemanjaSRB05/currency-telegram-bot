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
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализируем логгер
	if err := logger.InitGlobal(cfg.LogLevel); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Обработка команды миграций
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		if err := runMigrations(cfg); err != nil {
			logger.S.Fatalf("Migration failed: %v", err)
		}
		return
	}

	// Создаем и запускаем приложение
	application := app.New(cfg)
	if err := application.Run(); err != nil {
		logger.S.Errorf("Application failed: %v", err)
		os.Exit(1)
	}
}

// runMigrations выполняет миграции базы данных
func runMigrations(cfg *config.Config) error {
	db, err := goose.OpenDBWithDriver("postgres", cfg.DBURL)
	if err != nil {
		return err
	}
	defer db.Close()

	// Получаем команду миграции (по умолчанию - up)
	command := "up"
	if len(os.Args) > 2 {
		command = os.Args[2]
	}

	// Выполняем миграцию
	if err := goose.RunContext(context.Background(), command, db, "migrations"); err != nil {
		return err
	}

	logger.S.Info("Migrations completed successfully")
	return nil
}
