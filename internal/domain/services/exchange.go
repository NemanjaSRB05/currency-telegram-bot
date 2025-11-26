package services

import (
	"context"
)

// ExchangeService определяет операции для работы с курсами валют
type ExchangeService interface {
	GetRate(ctx context.Context, from, to string) (float64, error)
	ConvertAmount(ctx context.Context, amount float64, from, to string) (float64, error)
	GetSupportedCurrencies(ctx context.Context) (map[string]string, error) // код -> название
}

// ExchangeProvider определяет контракт для провайдеров курсов валют
type ExchangeProvider interface {
	GetRate(ctx context.Context, from, to string) (float64, error)
	GetName() string
	IsAvailable() bool
}
