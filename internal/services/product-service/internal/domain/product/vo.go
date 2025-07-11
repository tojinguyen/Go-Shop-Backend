package domain

import (
	"errors"
)

type Price struct {
	amount   float64
	currency string
}

func NewPrice(amount float64, currency string) (Price, error) {
	if amount < 0 {
		return Price{}, errors.New("price amount cannot be negative")
	}
	if currency == "" {
		return Price{}, errors.New("price currency cannot be empty")
	}
	return Price{amount: amount, currency: currency}, nil
}

func (p Price) GetAmount() float64 {
	return p.amount
}

func (p Price) GetCurrency() string {
	return p.currency
}
