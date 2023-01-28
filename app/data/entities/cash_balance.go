package entities

import "github.com/google/uuid"

type CashBalance struct {
	EntityBase
	PortfolioID   uuid.UUID
	CurrencyCode  string
	AmountInCents int64
}

type CashBalances []*CashBalance
