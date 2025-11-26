package cache

import (
	"sync"
	"time"
)

type cachedRate struct {
	rate      float64
	expiresAt time.Time
}

type RatesCache struct {
	mu    sync.RWMutex
	rates map[string]cachedRate // ключ: "USD_EUR", значение: кэшированный курс
	ttl   time.Duration
}

func NewRatesCache(ttlMinutes int) *RatesCache {
	return &RatesCache{
		rates: make(map[string]cachedRate),
		ttl:   time.Duration(ttlMinutes) * time.Minute,
	}
}

// Get возвращает курс из кэша, если он еще актуален
func (c *RatesCache) Get(from, to string) (float64, bool) {
	key := c.buildKey(from, to)

	c.mu.RLock()
	defer c.mu.RUnlock()

	cached, exists := c.rates[key]
	if !exists {
		return 0, false
	}

	if time.Now().After(cached.expiresAt) {
		return 0, false
	}

	return cached.rate, true
}

// Set сохраняет курс в кэш с TTL
func (c *RatesCache) Set(from, to string, rate float64) {
	key := c.buildKey(from, to)

	c.mu.Lock()
	defer c.mu.Unlock()

	c.rates[key] = cachedRate{
		rate:      rate,
		expiresAt: time.Now().Add(c.ttl),
	}
}

// buildKey создает ключ для кэша из пары валют
func (c *RatesCache) buildKey(from, to string) string {
	return from + "_" + to
}

// Cleanup удаляет просроченные записи (можно вызывать периодически)
func (c *RatesCache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, cached := range c.rates {
		if now.After(cached.expiresAt) {
			delete(c.rates, key)
		}
	}
}
