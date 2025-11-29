package cache

import (
	"testing"
)

func TestRatesCache(t *testing.T) {
	cache := NewRatesCache(1) // TTL 1 minute

	// Test Set and Get
	cache.Set("USD", "EUR", 0.85)
	rate, found := cache.Get("USD", "EUR")

	if !found {
		t.Error("Expected to find rate in cache")
	}
	if rate != 0.85 {
		t.Errorf("Expected rate 0.85, got %f", rate)
	}

	// Test non-existent rate
	_, found = cache.Get("USD", "GBP")
	if found {
		t.Error("Expected not to find non-existent rate")
	}
}

func TestRatesCacheExpiration(t *testing.T) {
	cache := NewRatesCache(1) // TTL 1 minute

	cache.Set("USD", "EUR", 0.85)

	// Rate should be available immediately
	_, found := cache.Get("USD", "EUR")
	if !found {
		t.Error("Rate should be available before expiration")
	}
}
