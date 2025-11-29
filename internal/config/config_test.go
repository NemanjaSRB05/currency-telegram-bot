package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestLoadConfig(t *testing.T) {
	// Сохраняем оригинальные значения
	originalBotToken := os.Getenv("BOT_TOKEN")
	originalDBURL := os.Getenv("DB_URL")

	// Устанавливаем тестовые значения
	os.Setenv("BOT_TOKEN", "test_token")
	os.Setenv("DB_URL", "postgres://test:test@localhost/test")

	// Сбрасываем Viper чтобы он перечитал environment variables
	viper.Reset()

	config, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if config.BotToken != "test_token" {
		t.Errorf("Expected BOT_TOKEN 'test_token', got '%s'", config.BotToken)
	}

	if config.DBURL != "postgres://test:test@localhost/test" {
		t.Errorf("Expected DB_URL 'postgres://test:test@localhost/test', got '%s'", config.DBURL)
	}

	// Восстанавливаем оригинальные значения
	if originalBotToken != "" {
		os.Setenv("BOT_TOKEN", originalBotToken)
	} else {
		os.Unsetenv("BOT_TOKEN")
	}
	if originalDBURL != "" {
		os.Setenv("DB_URL", originalDBURL)
	} else {
		os.Unsetenv("DB_URL")
	}
}

func TestLoadConfigMissingRequired(t *testing.T) {
	// Сохраняем оригинальные значения
	originalBotToken := os.Getenv("BOT_TOKEN")
	originalDBURL := os.Getenv("DB_URL")

	// Убираем обязательные переменные
	os.Unsetenv("BOT_TOKEN")
	os.Unsetenv("DB_URL")

	// Сбрасываем Viper чтобы он перечитал environment variables
	viper.Reset()

	_, err := Load()
	if err == nil {
		t.Error("Expected error for missing required environment variables")
	}

	// Восстанавливаем оригинальные значения
	if originalBotToken != "" {
		os.Setenv("BOT_TOKEN", originalBotToken)
	}
	if originalDBURL != "" {
		os.Setenv("DB_URL", originalDBURL)
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	// Сохраняем оригинальные значения
	originalBotToken := os.Getenv("BOT_TOKEN")
	originalDBURL := os.Getenv("DB_URL")
	originalLogLevel := os.Getenv("LOG_LEVEL")

	// Устанавливаем только обязательные значения
	os.Setenv("BOT_TOKEN", "test_token")
	os.Setenv("DB_URL", "postgres://test:test@localhost/test")
	os.Unsetenv("LOG_LEVEL")

	// Сбрасываем Viper чтобы он перечитал environment variables
	viper.Reset()

	config, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Проверяем значения по умолчанию
	if config.LogLevel != "info" {
		t.Errorf("Expected default LOG_LEVEL 'info', got '%s'", config.LogLevel)
	}

	if config.CacheTTLMinutes != 5 {
		t.Errorf("Expected default CACHE_TTL_MINUTES 5, got %d", config.CacheTTLMinutes)
	}

	// Восстанавливаем оригинальные значения
	if originalBotToken != "" {
		os.Setenv("BOT_TOKEN", originalBotToken)
	} else {
		os.Unsetenv("BOT_TOKEN")
	}
	if originalDBURL != "" {
		os.Setenv("DB_URL", originalDBURL)
	} else {
		os.Unsetenv("DB_URL")
	}
	if originalLogLevel != "" {
		os.Setenv("LOG_LEVEL", originalLogLevel)
	}
}
