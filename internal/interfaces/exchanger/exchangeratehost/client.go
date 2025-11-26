package exchangeratehost

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ExchangeRateHostResponse struct {
	Result string             `json:"result"`
	Rates  map[string]float64 `json:"rates"`
}

type ExchangeRateHostClient struct {
	baseURL    string
	httpClient *http.Client
}

func New() *ExchangeRateHostClient {
	return &ExchangeRateHostClient{
		baseURL: "https://api.exchangerate.host",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetRate получает курс обмена между двумя валютами
func (c *ExchangeRateHostClient) GetRate(ctx context.Context, from, to string) (float64, error) {
	url := fmt.Sprintf("%s/latest?base=%s&symbols=%s", c.baseURL, from, to)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	var apiResponse ExchangeRateHostResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return 0, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if apiResponse.Result != "success" {
		return 0, fmt.Errorf("API returned error result: %s", apiResponse.Result)
	}

	rate, exists := apiResponse.Rates[to]
	if !exists {
		return 0, fmt.Errorf("rate for %s not found in response", to)
	}

	return rate, nil
}

// GetName возвращает имя провайдера
func (c *ExchangeRateHostClient) GetName() string {
	return "exchangerate.host"
}

// IsAvailable проверяет доступность провайдера
func (c *ExchangeRateHostClient) IsAvailable() bool {
	// Простая проверка доступности
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.GetRate(ctx, "USD", "EUR")
	return err == nil
}
