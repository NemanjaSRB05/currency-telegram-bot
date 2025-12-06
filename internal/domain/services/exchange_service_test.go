package services_test

import (
	"context"
	"testing"

	"github.com/crocxdued/currency-telegram-bot/internal/domain/services"
	"github.com/crocxdued/currency-telegram-bot/internal/interfaces/repository/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockExchangeProvider struct {
	mock.Mock
}

func (m *MockExchangeProvider) GetRate(ctx context.Context, from, to string) (float64, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockExchangeProvider) GetName() string {
	return "mock-provider"
}

func (m *MockExchangeProvider) IsAvailable() bool {
	// ВАЖНО: используем мок, а не фиксированное значение
	args := m.Called()
	return args.Bool(0)
}

func TestExchangeService_GetRate(t *testing.T) {
	mockProvider := &MockExchangeProvider{}
	cache := cache.NewRatesCache(5)

	mockProvider.On("IsAvailable").Return(true)
	mockProvider.On("GetRate", mock.Anything, "USD", "EUR").Return(0.85, nil)

	service := services.NewExchangeService([]services.ExchangeProvider{mockProvider}, cache)

	rate, err := service.GetRate(context.Background(), "USD", "EUR")

	assert.NoError(t, err)
	assert.InDelta(t, 0.85, rate, 0.001)
	mockProvider.AssertExpectations(t)
}

func TestExchangeService_ConvertAmount(t *testing.T) {
	mockProvider := &MockExchangeProvider{}
	cache := cache.NewRatesCache(5)

	mockProvider.On("IsAvailable").Return(true)
	mockProvider.On("GetRate", mock.Anything, "USD", "EUR").Return(0.85, nil)

	service := services.NewExchangeService([]services.ExchangeProvider{mockProvider}, cache)

	result, err := service.ConvertAmount(context.Background(), 100.0, "USD", "EUR")

	assert.NoError(t, err)
	assert.InDelta(t, 85.0, result, 0.001)
	mockProvider.AssertExpectations(t)
}

func TestExchangeService_NoAvailableProvider(t *testing.T) {
	mockProvider := &MockExchangeProvider{}
	cache := cache.NewRatesCache(5)

	// Провайдер недоступен — GetRate НЕ должен быть вызван
	mockProvider.On("IsAvailable").Return(false)

	service := services.NewExchangeService([]services.ExchangeProvider{mockProvider}, cache)

	_, err := service.GetRate(context.Background(), "USD", "EUR")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get exchange rate")
	mockProvider.AssertNotCalled(t, "GetRate", mock.Anything, mock.Anything, mock.Anything)
	mockProvider.AssertExpectations(t)
}
