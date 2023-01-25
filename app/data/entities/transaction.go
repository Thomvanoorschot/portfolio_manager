package entities

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type TransactionType int

const (
	Unknown TransactionType = iota + 1
	Buy
	Sell
	Deposit
	Withdrawal
)

func ConvertToTransactionType(transactionType string) (TransactionType, error) {
	switch transactionType {
	case "Koop":
		return Buy, nil
	case "Verkoop":
		return Sell, nil
	default:
		return Unknown, errors.New("could not convert transactionType")
	}
}

type Transaction struct {
	TransactedAt      time.Time
	CurrencyCode      string
	TransactionType   TransactionType
	Product           string
	ISIN              string
	Description       string
	Amount            float64
	PriceInCents      int64
	CommissionInCents int64
	ExternalId        string
	PortfolioID       uuid.UUID
	Symbol            string
	UniqueHash        string
	EntityBase
}

type Transactions []*Transaction
