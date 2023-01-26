package entities

import (
	"errors"
	"github.com/Thomvanoorschot/portfolioManager/app/enums"
	"github.com/google/uuid"
	"time"
)

func ConvertToTransactionType(transactionType string) (enums.TransactionType, error) {
	switch transactionType {
	case "Koop":
		return enums.Buy, nil
	case "Verkoop":
		return enums.Sell, nil
	default:
		return enums.Unknown, errors.New("could not convert transactionType")
	}
}

type Transaction struct {
	EntityBase
	TransactedAt      time.Time
	CurrencyCode      string
	TransactionType   enums.TransactionType
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
}

type Transactions []*Transaction
