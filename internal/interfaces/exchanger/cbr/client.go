package cbr

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	ID       string `xml:"ID,attr"`
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  string `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

type CBRClient struct {
	baseURL    string
	httpClient *http.Client
}

func New() *CBRClient {
	return &CBRClient{
		baseURL: "https://www.cbr.ru",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetRate получает курс обмена (все курсы от RUB)
func (c *CBRClient) GetRate(ctx context.Context, from, to string) (float64, error) {
	// CBR API предоставляет курсы только относительно RUB
	if from != "RUB" && to != "RUB" {
		return 0, fmt.Errorf("CBR API supports only RUB conversions")
	}

	url := fmt.Sprintf("%s/scripts/XML_daily_eng.asp", c.baseURL)

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

	var valCurs ValCurs
	if err := xml.Unmarshal(body, &valCurs); err != nil {
		return 0, fmt.Errorf("failed to parse XML: %w", err)
	}

	// Определяем целевую валюту
	targetCurrency := to
	if from == "RUB" {
		targetCurrency = to
	} else {
		targetCurrency = from
	}

	// Ищем нужную валюту
	var rate float64
	var found bool

	for _, valute := range valCurs.Valutes {
		if valute.CharCode == targetCurrency {
			// Конвертируем строку в float (заменяем запятую на точку)
			valueStr := strings.Replace(valute.Value, ",", ".", -1)
			rate, err = strconv.ParseFloat(valueStr, 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse rate: %w", err)
			}

			// Делим на номинал (например, 1 USD = 90 RUB, но может быть 10 CNY = 120 RUB)
			nominal, _ := strconv.ParseFloat(valute.Nominal, 64)
			rate = rate / nominal

			found = true
			break
		}
	}

	if !found {
		return 0, fmt.Errorf("currency %s not found in CBR response", targetCurrency)
	}

	// Если конвертируем из RUB в другую валюту, инвертируем курс
	if from == "RUB" {
		rate = 1 / rate
	}

	return rate, nil
}

// GetName возвращает имя провайдера
func (c *CBRClient) GetName() string {
	return "cbr.ru"
}

// IsAvailable проверяет доступность провайдера
func (c *CBRClient) IsAvailable() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.GetRate(ctx, "USD", "RUB")
	return err == nil
}
