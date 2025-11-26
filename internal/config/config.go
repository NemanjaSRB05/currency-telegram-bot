package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	BotToken        string `mapstructure:"BOT_TOKEN"`
	DBURL           string `mapstructure:"DB_URL"`
	LogLevel        string `mapstructure:"LOG_LEVEL"`
	CacheTTLMinutes int    `mapstructure:"CACHE_TTL_MINUTES"`
}

func Load() (*Config, error) {
	// Настраиваем Viper для чтения .env файлов
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Устанавливаем значения по умолчанию
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("CACHE_TTL_MINUTES", 5)

	// Читаем конфиг файл (игнорируем ошибку если файла нет)
	_ = viper.ReadInConfig()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Валидация обязательных полей
	if config.BotToken == "" {
		return nil, fmt.Errorf("BOT_TOKEN is required")
	}
	if config.DBURL == "" {
		return nil, fmt.Errorf("DB_URL is required")
	}

	return &config, nil
}
