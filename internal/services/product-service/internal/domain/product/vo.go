package domain

import "errors"

type Price struct {
	Amount   float64
	Currency string
}

func NewPrice(amount float64, currency string) (Price, error) {
	if amount < 0 {
		return Price{}, errors.New("price amount cannot be negative")
	}
	if currency == "" {
		return Price{}, errors.New("price currency cannot be empty")
	}
	return Price{Amount: amount, Currency: currency}, nil
}
