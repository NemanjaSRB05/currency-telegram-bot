package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/crocxdued/currency-telegram-bot/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type App struct {
	config *config.Config
	db     *sqlx.DB
}

func New(cfg *config.Config) *App {
	return &App{
		config: cfg,
	}
}

// initDB инициализирует подключение к базе данных
func (a *App) initDB(ctx context.Context) error {
	db, err := sqlx.ConnectContext(ctx, "postgres", a.config.DBURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	a.db = db
	log.Println("Database connection established")
	return nil
}

// Run запускает приложение
func (a *App) Run() error {
	ctx := context.Background()

	// Инициализируем БД
	if err := a.initDB(ctx); err != nil {
		return err
	}
	defer a.db.Close()

	// TODO: Инициализировать и запустить бота

	log.Println("Application started successfully")
	<-ctx.Done()
	return nil
}
