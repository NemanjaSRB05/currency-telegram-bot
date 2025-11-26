package entities

import (
	"time"
)

// Currency представляет валюту с кодом и названием
type Currency struct {
	Code string
	Name string
}

// ExchangeRate представляет курс обмена между двумя валютами
type ExchangeRate struct {
	From        Currency
	To          Currency
	Rate        float64
	LastUpdated time.Time
}

// UserFavorite представляет избранную пару валют пользователя
type UserFavorite struct {
	ID           int64
	UserID       int64
	FromCurrency string
	ToCurrency   string
	CreatedAt    time.Time
}
