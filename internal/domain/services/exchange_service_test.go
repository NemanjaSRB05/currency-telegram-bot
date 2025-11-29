package services

import (
	"context"
	"testing"

	"github.com/crocxdued/currency-telegram-bot/internal/interfaces/repository/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockExchangeProvider для тестов
type MockExchangeProvider struct {
	mock.Mock
}

func (m *MockExchangeProvider) GetRate(ctx context.Context, from, to string) (float64, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockExchangeProvider) GetName() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockExchangeProvider) IsAvailable() bool {
	args := m.Called()
	return args.Bool(0)
}

func TestExchangeService_GetRate(t *testing.T) {
	// Создаем моки
	mockProvider := new(MockExchangeProvider)
	cache := cache.NewRatesCache(5)

	// Настраиваем ожидания
	mockProvider.On("IsAvailable").Return(true)
	mockProvider.On("GetName").Return("test-provider")
	mockProvider.On("GetRate", mock.Anything, "USD", "EUR").Return(0.85, nil)

	// Создаем сервис
	service := NewExchangeService([]ExchangeProvider{mockProvider}, cache)

	// Тестируем
	rate, err := service.GetRate(context.Background(), "USD", "EUR")

	// Проверяем
	assert.NoError(t, err)
	assert.Equal(t, 0.85, rate)
	mockProvider.AssertCalled(t, "GetRate", mock.Anything, "USD", "EUR")
}

func TestExchangeService_ConvertAmount(t *testing.T) {
	mockProvider := new(MockExchangeProvider)
	cache := cache.NewRatesCache(5)

	mockProvider.On("IsAvailable").Return(true)
	mockProvider.On("GetName").Return("test-provider")
	mockProvider.On("GetRate", mock.Anything, "USD", "EUR").Return(0.85, nil)

	service := NewExchangeService([]ExchangeProvider{mockProvider}, cache)

	// Конвертируем 100 USD в EUR
	result, err := service.ConvertAmount(context.Background(), 100, "USD", "EUR")

	assert.NoError(t, err)
	assert.Equal(t, 85.0, result) // 100 * 0.85 = 85
}
